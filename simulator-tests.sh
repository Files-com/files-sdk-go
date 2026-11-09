#!/usr/bin/env bash
# Runs the Go SDK's tests that need the Files.com mock server's simulator, which `go test ./...` skips without one.
# Each group runs in one targets/mock-server/with-scenario.rb launch of its own, under its scenario, from the generated
# mock server named by the first argument, with the race detector and FILES_API_SCENARIOS=required, so a test that would
# be skipped fails. The failure-cases group runs the shared failure-case matrix (user.TestFailureCases), whose launcher
# adds the scenario's fault rules and fails unless every case passed with exactly its expected attempts and decisions.
# go test's JSON events must then show every listed test run once and pass once, with no skip or failure anywhere in
# the run, each listed matrix with exactly its cases passed, and the package passed. The launcher's records and the
# go test events of each group go to <output directory>/<group>, which must not exist yet.
#
#   ./simulator-tests.sh ../mock-server "$OUT"
set -euo pipefail
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
mock="${1:?generated mock server directory required}"
out="${2:?output directory required}"

# The file tests share one simulator, as the launcher's api-scenarios-canonical reset left it: none resets it, and each
# works under remote paths of its own. TestDownload and TestUploader run lib.PathSpec's 12 cases under one group each.
file_tests=(
  TestFS_Open TestClient_UploadFile
  TestClient_UploadFolder TestClient_UploadFolder_Dot TestClient_UploadFolder_Relative TestClient_Uploader
  TestClient_UploadFile_To_Existing_Dir TestClient_UploadFile_To_NonExistingPath
  TestClient_UploadFile_To_NonExistingPath_WithSlash TestClient_UploadFolder_as_file2
  TestClient_UploadFolder_RemotePathWithStartingSlash TestClient_UploadFolder_ZeroByteFile
  TestClient_DownloadFolder TestClient_DownloadFolder_Smart TestClient_DownloadFolder_file_to_file
  TestClient_DownloadFolder_file_to_implicit TestClient_DownloadFolder_file_only
  TestClient_Downloader_Delete_Source TestClient_Downloader_Delete_Source_Folder TestClient_Downloader_Move_Source
  TestClient_UploadFolder_Move_Source TestClient_UploadFolder_Move_Source_Missing_Dir
  TestClient_Downloader_Move_Source_Missing_Dir TestClient_UploadFile_Delete_Source TestClient_Uploader_Files
  TestClient_Uploader_File_Remote_Folder TestClient_Uploader_File_Remote_Path_Error TestClient_Uploader_Directories
  TestClient_Uploader_Empty_Directory TestClient_Uploader_Directories_With_DeleteSource TestClient_ListFor
  TestDownload TestUploader
)
file_matrices="TestDownload/downloader=12 TestUploader/uploader=12"
# TestClient_Wait resets the simulator to the profile it needs, so it runs in a launch of its own.
filemigration_tests=(TestClient_Wait)

run_group() {
  local group="$1" package="$2" matrices="$3" scenario="$4"
  shift 4
  local run="${out}/${group}"
  [[ ! -e "${run}" ]] || { echo "${run} already exists" >&2; exit 1; }
  mkdir -p "${run}"
  local selector
  selector="^($(IFS='|'; echo "$*"))\$"
  local status=0
  (
    cd "${DIR}"
    FILES_API_SCENARIOS=required ruby "${mock}/with-scenario.rb" --scenario ${scenario} --out "${run}/launcher" \
      --startup-timeout 30 --consumer-timeout 180 --settle-timeout 5 -- \
      bash -c 'go test -race -count=1 -json -run "$1" "$2" > "$3"' consumer "${selector}" "${package}" "${run}/go-test.jsonl"
  ) || status=$?
  ruby -rjson -e '
    log, package, matrices, *names = ARGV
    abort "#{log} was not written: go test did not run" unless File.exist?(log)
    events = File.readlines(log, chomp: true).filter_map { |line| JSON.parse(line) if line.start_with?("{") }
    failures = []
    events.select { |event| %w[fail skip].include?(event["Action"]) }.each do |event|
      failures << "#{event["Test"] || event["Package"]} reported #{event["Action"]}"
    end
    names.each do |name|
      actions = events.select { |event| event["Test"] == name }.map { |event| event["Action"] }
      failures << "#{name}: #{actions.inspect}, not one run and one pass" unless actions.count("run") == 1 && actions.count("pass") == 1
    end
    matrices.split.each do |matrix|
      prefix, count = matrix.split("=")
      passed = events.select { |event| event["Action"] == "pass" && event["Test"].to_s.start_with?("#{prefix}/") }.map { |event| event["Test"] }
      failures << "#{prefix}: #{passed.size} cases passed, not #{count}" unless passed.uniq.size == passed.size && passed.size == Integer(count)
    end
    results = events.select { |event| event["Test"].nil? && event["Package"].to_s.end_with?("/#{package.delete_prefix("./")}") }.map { |event| event["Action"] }
    failures << "package #{package}: #{results.inspect}, not a pass" unless results.include?("pass")
    failed = events.select { |event| event["Action"] == "fail" }.map { |event| event["Test"] }.compact
    events.each { |event| $stderr.print event["Output"] if event["Action"] == "output" && failed.include?(event["Test"]) }
    abort "#{log}:\n  #{failures.join("\n  ")}" unless failures.empty?
    puts "#{log}: #{names.size} tests passed, each once#{matrices.empty? ? "" : ", with #{matrices}"}"
  ' "${run}/go-test.jsonl" "${package}" "${matrices}" "$@" || status=1
  return "${status}"
}

status=0
run_group file ./file "${file_matrices}" api-scenarios-canonical "${file_tests[@]}" || status=1
run_group filemigration ./filemigration "" api-scenarios-canonical "${filemigration_tests[@]}" || status=1
# The SDK's own retries resend a dropped request (the scenario's sdk-network-retries variant).
run_group failure-cases ./user "" "failure-cases --failure-variant sdk-network-retries" TestFailureCases || status=1
exit "${status}"
