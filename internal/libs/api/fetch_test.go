package api

import (
	"encoding/json"
	"testing"
	"time"
)

func TestFetch(t *testing.T) {
	l, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Error(err)
	}
	url := CreateURL(URLOptions{
		StartRecord:    2,
		MaximumRecords: 1,
		From:           time.Date(2022, 3, 1, 0, 0, 0, 0, l),
		Until:          time.Date(2022, 5, 1, 0, 0, 0, 0, l),
		Any:            "防災",
	})

	exp := `{"numberOfRecords":582,"numberOfReturn":1,"startRecord":2,"nextRecordPosition":3,"speechRecord":[{"speechID":"120804194X00320220428_008","issueID":"120804194X00320220428","imageKind":"会議録","searchObject":8,"session":208,"nameOfHouse":"衆議院","nameOfMeeting":"原子力問題調査特別委員会","issue":"第3号","date":"2022-04-28","closing":null,"speechOrder":8,"speaker":"井林辰憲","speakerYomi":"いばやしたつのり","speakerGroup":"自由民主党","speakerPosition":null,"speakerRole":null,"speech":"○井林委員　ありがとうございます。\r\n　やはり最後は組織の厚みというのが重要になるのかなというふうに思いますが、ここは二律背反なんですが、組織の厚みというのはその産業界の厚みにつながりますので、そうすると、これは推進と規制の議論からするとまた難しい問題を抱えているのかなというふうに思っております。今日は、この場ではこれ以上は踏み込むのはやめておきたいと思います。\r\n　あと、この場で個別の原子力サイトについてやり取りをするのは不適切だと思うんですが、私の選挙区外ということと、このサイトということではないんですが、審査の方向性としてあるべき方向に向かっているし、これは一度、委員長にちゃんと確認をさせていただきたいと思うので、取り上げさせていただきたいと思います。\r\n　泊原発の三号機でございます。私が内閣府の原子力防災担当の政務官を務めさせていただいた二〇一六年に総合防災訓練も行われました。ですので、折につけ、審査会合など、毎回見ているわけではないですが、ウォッチをしてまいりました。\r\n　正直言って、規制委員会と事業者とのやり取りがうまくできていないなというのが、もう少しかみ合った議論をすれば安全性も高まるし、これは私の地元でもそうなんですが、やはりそういう議論が地元の皆様の安心というものにつながっていくんじゃないかなというふうに思って見ておりました。\r\n　その中で、先月の、三月三十一日の審査会合で、審査における残された論点の確認というのが行われました。これは、規制当局側から、気になる点というか考えをしっかりと示したということ。そして、その後の四月十二日の原子力規制委員会と泊原発の事業者である北海道電力の経営層、これは議事録が手元にありますが、社長も出席をされております。\r\n　この場で意見交換が行われまして、ちょっとこの議事録をそのまま読み上げさせていただきますと、更田委員長の発言ですが、スペシフィックな評価や解析に関する担当者同士の審査会合というようなもの、我々の方としては応じることができると思っている。この前段では、もちろん、原則公開とか、ユーチューブでちゃんと上げるということです。\r\n　これは、翻訳をすると、私は、審査会合というのは今まで、規制委員会の委員の方々が出てきて行うというのが審査会合だということですが、ただ、規制委員会の委員が判断すべきことと、それを支える事務組織である規制庁の職員が判断すべきことというのは、それぞれいろいろあるんじゃないかなということを思っておりました。\r\n　そういうことも踏まえると、私、これはちょっと前向きなというか、私のそういう思いも込めて判断をすると、審査会合は審査会合でやるんだけれども、事務的なヒアリングとは別に、公開を原則としながら、事務的な審査会合というんですかね、そういうものも考えられるんだというような発言だというふうに思っておりますし、これは私は、今回のことをロールモデルとして、イメージが勝手に一緒だと解釈するのも危険なので、ここをちゃんと確認をさせていただきたいと思いますが、今回の泊の三号機における一連の審査及び発言について、今後の方向性と、他のサイトの審査への展開について、委員長の所見をというか考えをお伺いをしたいと思います。","startPage":2,"speechURL":"https://kokkai.ndl.go.jp/txt/120804194X00320220428/8","meetingURL":"https://kokkai.ndl.go.jp/txt/120804194X00320220428","pdfURL":"https://kokkai.ndl.go.jp/img/120804194X00320220428/2"}]}`
	var expJson *SpeechJson
	if err := json.Unmarshal([]byte(exp), &expJson); err != nil {
		t.Error(err)
	}

	result := Fetch(url)
	if result.Err != nil {
		t.Error(result.Err)
	}

	if len(result.SpeechJson.SpeechRecord) != len(expJson.SpeechRecord) {
		t.Errorf("Expected len(resp) is %d, but got='%d'", len(expJson.SpeechRecord), len(result.SpeechJson.SpeechRecord))
	}
	if result.SpeechJson.SpeechRecord[0].Speech != expJson.SpeechRecord[0].Speech {
		t.Errorf("Expected speech differs from actual")
	}

}
