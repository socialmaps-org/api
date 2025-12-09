package moderation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCleanup(t *testing.T) {
	const exp = "This is utterly unbelievable! 🤬"
	cases := []string{
		"𝔗𝔥𝔦𝔰 𝔦𝔰 𝔲𝔱𝔱𝔢𝔯𝔩𝔶 𝔲𝔫𝔟𝔢𝔩𝔦𝔢𝔳𝔞𝔟𝔩𝔢! 🤬",
		"𝕿𝖍𝖎𝖘 𝖎𝖘 𝖚𝖙𝖙𝖊𝖗𝖑𝖞 𝖚𝖓𝖇𝖊𝖑𝖎𝖊𝖛𝖆𝖇𝖑𝖊! 🤬",
		"𝓣𝓱𝓲𝓼 𝓲𝓼 𝓾𝓽𝓽𝓮𝓻𝓵𝔂 𝓾𝓷𝓫𝓮𝓵𝓲𝓮𝓿𝓪𝓫𝓵𝓮! 🤬",
		"𝒯𝒽𝒾𝓈 𝒾𝓈 𝓊𝓉𝓉𝑒𝓇𝓁𝓎 𝓊𝓃𝒷𝑒𝓁𝒾𝑒𝓋𝒶𝒷𝓁𝑒! 🤬",
		"𝕋𝕙𝕚𝕤 𝕚𝕤 𝕦𝕥𝕥𝕖𝕣𝕝𝕪 𝕦𝕟𝕓𝕖𝕝𝕚𝕖𝕧𝕒𝕓𝕝𝕖! 🤬",
		"𝐓𝐡𝐢𝐬 𝐢𝐬 𝐮𝐭𝐭𝐞𝐫𝐥𝐲 𝐮𝐧𝐛𝐞𝐥𝐢𝐞𝐯𝐚𝐛𝐥𝐞! 🤬",
		"𝗧𝗵𝗶𝘀 𝗶𝘀 𝘂𝘁𝘁𝗲𝗿𝗹𝘆 𝘂𝗻𝗯𝗲𝗹𝗶𝗲𝘃𝗮𝗯𝗹𝗲! 🤬",
		"𝘛𝘩𝘪𝘴 𝘪𝘴 𝘶𝘵𝘵𝘦𝘳𝘭𝘺 𝘶𝘯𝘣𝘦𝘭𝘪𝘦𝘷𝘢𝘣𝘭𝘦! 🤬",
		"𝙏𝙝𝙞𝙨 𝙞𝙨 𝙪𝙩𝙩𝙚𝙧𝙡𝙮 𝙪𝙣𝙗𝙚𝙡𝙞𝙚𝙫𝙖𝙗𝙡𝙚! 🤬",
	}

	for _, tc := range cases {
		act := CleanUp(tc)

		require.Equal(t, exp, act, tc)
	}
}
