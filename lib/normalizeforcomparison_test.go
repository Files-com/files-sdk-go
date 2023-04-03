package lib

import (
	"testing"
)

func TestNormalizeForComparison(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"filename.txt", "filename.txt"},
		{"FiLeNaMe.TxT", "filename.txt"},
		{"FILENAME.TXT", "filename.txt"},
		{"FÎŁĘÑÂMÉ.TXT", "filename.txt"},
		{"Fïłèńämê.Txt", "filename.txt"},
		{"a/b/c.txt", "a/b/c.txt"},
		{"A\\B\\C.TXT", "a/b/c.txt"},
		{"A/B\\C.TXT", "a/b/c.txt"},
		{"//a/b//c.txt", "a/b/c.txt"},
		{"a/b/c.txt  ", "a/b/c.txt"},
		{"a/b/c.txt\t", "a/b/c.txt"},
		{"a/b/c.txt\n", "a/b/c.txt"},
		{"a/b/c.txt\r", "a/b/c.txt"},
		{" space_at_beginning", " space_at_beginning"},
		{"space_at_end ", "space_at_end"},
		{"tab\tseperated", "tab\tseperated"},
		{"<title>hello</hello>", "<title>hello</hello>"},
		{"안녕하세요", "안녕하세요"},
		{"こんにちは", "こんにちは"},
		{"今日は", "今日は"},
		{"longest_unicode_character_﷽", "longest_unicode_character_﷽"},
		{"invalid_null_byte_before\u0000after", "invalid_null_byte_beforeafter"},
		{"a/b/c/../../hello", "a/b/c/hello"},
		{"a/b/c/././hello", "a/b/c/hello"},
		{"one_code_point_ą", "one_code_point_a"},
		{"two_code_points_ą", "two_code_points_a"},
		{"one_code_point_훯", "one_code_point_훯"},
		{"three_code_points_훯", "three_code_points_훯"},
		{"ÞþŊŋŦŧ", "þþŋŋŧŧ"},
		{"ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖØÙÚÛÜÝßàáâãäåæçèéêëìíîïðñòóôõöøùúûüýÿ", "aaaaaaaeceeeeiiiidnoooooouuuuyssaaaaaaaeceeeeiiiidnoooooouuuuyy"},
		{"ĀāĂăĄąĆćĈĉĊċČčĎďĐđĒēĔĕĖėĘęĚěĜĝĞğĠġĢģĤĥĦħĨĩĪīĬĭĮįİĲĳ", "aaaaaaccccccccddddeeeeeeeeeegggggggghhhhiiiiiiiiiijij"},
		{"ĴĵĶķĹĺĻļĽľŁłŃńŅņŇňŉŌōŎŏŐőŒœŔŕŖŗŘřŚśŜŝŞşŠšŢţŤť", "jjkkllllllllnnnnnnʼnoooooooeoerrrrrrsssssssstttt"},
		{"ŨũŪūŬŭŮůŰűŲųŴŵŶŷŸŹźŻżŽž", "uuuuuuuuuuuuwwyyyzzzzzz"},
		{"😂❤️😍🤣😊🙏💕😭😘👍😅👏😁♥️🔥💔💖💙😢🤔😆🙄💪😉☺️👌🤗", "😂❤️😍🤣😊🙏💕😭😘👍😅👏😁♥️🔥💔💖💙😢🤔😆🙄💪😉☺️👌🤗"},
		{"💜😔😎😇🌹🤦🎉💞✌️✨🤷😱😌🌸🙌😋💗💚😏💛🙂💓🤩😄😀🖤😃💯🙈👇🎶😒🤭❣️", "💜😔😎😇🌹🤦🎉💞✌️✨🤷😱😌🌸🙌😋💗💚😏💛🙂💓🤩😄😀🖤😃💯🙈👇🎶😒🤭❣️"},
		{"emoji_‼️", "emoji_!!️"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			output := NormalizeForComparison(tc.input)
			if output != tc.expected {
				t.Errorf("Expected %s but got %s", tc.expected, output)
			}
		})
	}
}
