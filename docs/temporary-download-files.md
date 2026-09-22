# Temporary download files

A folder download writes each file to a temporary path first and renames it into
place once it is complete. If a transfer is interrupted, a later run finds that
temporary file again and continues from the bytes already on disk.

## Names

Temporary names live in a namespace reserved for this client. Beside the file
they belong to, two forms are generated:

| Form | Name | When |
| --- | --- | --- |
| plain | `.~files-cli.<name>.download` | the usual case |
| encoded | `.~files-cli~<n>.<digest of the name>.<shortened name>.download` | the name is too long to fit in a path element, or another temporary download already holds the plain name |

`<digest of the name>` identifies the file the encoded name belongs to, so the
shortened name next to it is only there to be read. The character right after
`.~files-cli` — `.` or `~` — is what tells the two forms apart, and it is not
part of any file name, so no file name can be spelled to produce another file's
temporary name.

Inside an external temporary directory (`TempPath`), files from every folder of
a download, and from every download that shares the directory, are staged side
by side, so a name alone does not say which file a temporary download belongs
to. There, a temporary download is always encoded and identified by the file's
whole local path — its destination directory and the path below it, made
absolute — instead of by its name:

| Form | Name | When |
| --- | --- | --- |
| external | `.~files-cli~<n>.path-<digest of the full path>.<shortened name>.download` | always, inside `TempPath` |

The `path-` marker is not a hexadecimal digit, so an external name is never
spelled like a name identified by a name digest. Two files that differ anywhere
in their paths — `a/x.bin` and `b/x.bin` under one destination, or the same
`a/x.bin` under two destination roots — have two temporary downloads, and one
file reached through different spellings of its path (an explicit output path,
a destination root plus the path below it, or a path relative to the working
directory) has one.

On macOS a temporary download is a folder with that name, holding the file under
its own name. On Linux and Windows it is a file with that name.

The short-lived names used while moving a finished file between directories are
in the same namespace.

## Known limitation: an adaptive download stopped by a process kill

A later run continues a temporary download from the bytes already in it, and
finishes one that is already the size of the file without asking the server
for content. An ordinary download appends in order, so its temporary file is
always a valid prefix. An adaptive download preallocates the whole file and
writes ranges at their final offsets; it cuts the file back to the valid prefix
when it is interrupted cleanly, but a process that is killed cannot. The
full-size file it leaves is then judged by its size like any other, so a later
run may finish it without fetching the ranges that were never written. This
behavior predates the reserved naming and is unchanged by it. Delete the
temporary file of a killed adaptive download before running it again.

## Publishing

A download is written to its temporary path through one open file, and the
transfer remembers which file that is. When the download is complete it is not
renamed into place directly. It is first *claimed*: moved to a random name that
is already a valid 8.3 name, so NTFS keeps no second name for it, inside a
folder in the reserved namespace. When the temporary file already sits in such
a folder, as a macOS temporary download usually does, that folder is used;
otherwise — a file on Linux or Windows, or a temporary file the caller named
explicitly through a checkpoint — a fresh `.~files-cli~<random>.download`
folder is created beside it. The file found under the claimed name is compared with the file
the transfer wrote. Only if they are the same file is it moved to its final
name; otherwise the download fails with an error saying the temporary file was
replaced, and nothing is published. The claim folder is removed once the file
has left it.

This is what stops a file that reaches the temporary path under another name
the filesystem gives it — such as the `~FILES~1.DOW` that NTFS derives from a
long name — from being published as this download. On Windows the transfer
check also resolves a local path to the name the filesystem stores, so a name
that is only another spelling of a temporary download is not transferred, and a
path that exists but cannot be resolved is refused rather than allowed.

If the move to the final name fails, for example because the job was paused,
the file is put back under its temporary name so a later run continues from it.
If the process stops between the claim and the move, the finished file is left
under its claimed name inside the reserved folder: it is never transferred, but
it is not found by a later run either, which downloads the file again. Delete
such leftovers when they are no longer wanted.

## These names never transfer

A path whose name — or whose parent folder's name — is in this namespace is never
uploaded and never downloaded, in either direction, whatever `Ignore` and
`Include` rules the caller passes. Both the local path and the remote path of a
transfer are checked, because a folder selected as the root of a transfer loses
its own name on the way to the other side.

The comparison folds case the way `strings.EqualFold` does, because a filesystem
that finds a name whatever its case would otherwise let a differently spelled
name reach the same temporary path. That includes spellings that take a different
number of bytes: on macOS `.~fileſ-cli.X.download` (`ſ`, U+017F) and
`.~files-cli.X.download` are one and the same file.

Without this, a remote file named after another file's temporary download would
be written to that temporary path and then delivered as the file it was named
after, instead of that file's own content.

Names that only look similar are not in the namespace: `report.csv.download`,
`.~files-cli-notes.download` and `.~files-clients.download` are ordinary files,
so whether they transfer is decided by the caller's own `Ignore` and `Include`
rules as usual. The default rules, which apply when a caller passes none,
already exclude everything ending in `.download`.

## Upgrading from the previous naming

Before this naming, temporary downloads were `<name>.download`. A leftover file
with that name is now an ordinary file:

* It is not continued from. An interrupted download from an earlier version
  starts again from the beginning. This is deliberate — a genuine leftover from
  an earlier version cannot be told apart from a file planted under the same
  name, which is the problem this naming solves.
* It is not deleted or written to, and it is not in the reserved namespace, so
  it will be uploaded or downloaded like any other file. Delete leftovers when
  they are no longer wanted.

A temporary path recorded in a job checkpoint (`TmpPath`) is still used as given,
so an explicitly paused transfer resumes as it did before.

## Upgrading from name-only external temporary downloads

Before external temporary downloads were identified by the file's whole path, a
temporary download inside `TempPath` carried only the file's name:
`.~files-cli.<name>.download`, or the encoded form with a digest of the name.
Such a leftover does not say which of the files named that way it was for —
`a/x.bin` or `b/x.bin`, or the same path under another destination root — so
it is not attributed to any of them:

* It is not continued from. An interrupted download through `TempPath` from an
  earlier version starts again from the beginning, under a new temporary name
  beside the leftover. This is deliberate: continuing from it is what let one
  file be delivered with another file's bytes.
* It is not deleted or written to. It stays in the reserved namespace, so it is
  never transferred, but it is not cleaned up either. Delete leftovers in
  `TempPath` when they are no longer wanted.

A temporary path recorded in a job checkpoint (`TmpPath`) names the file
explicitly and is still used as given, so an explicitly paused transfer through
`TempPath` resumes as it did before.
