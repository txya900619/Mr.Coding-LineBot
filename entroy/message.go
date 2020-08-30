package entroy

import (
	"fmt"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
)

func StartMessage() *linebot.FlexMessage {
	flexContainer := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Size: linebot.FlexBubbleSizeTypeMega,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Size:   linebot.FlexTextSizeTypeXxl,
					Weight: linebot.FlexTextWeightTypeBold,
					Text:   "歡迎",
				},
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Size: linebot.FlexTextSizeTypeSm,
					Text: "在遊戲開始之前要告訴你一些資訊",
				},
			},
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Size: linebot.FlexTextSizeTypeXl,
					Text: "社團茶會",
				},
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: "時間：9 月 22 號",
				},
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: "地點：綜合科館 B1 第三演講廳",
				},
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: "期待你的參與歐～～～",
				},
			},
		},
		Footer: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type: linebot.FlexComponentTypeButton,
					Action: &linebot.URIAction{
						Label: "點我填寫報名表單",
						URI:   "https://docs.google.com/forms/d/1X1ZqMvUsA4cp4uPQKzH9iwc7dzNTEp2z3wjQzsv40gQ/viewform?edit_requested=true",
					},
				},
			},
		},
	}

	return linebot.NewFlexMessage("StartGameMessage", flexContainer)
}

func (p *Player) GetQuestionMessageByID(questionID uint) *linebot.FlexMessage {
	var title string
	answeredCount := len(p.AnsweredList)
	switch answeredCount {
	case 0:
		title = "第一題"
	case 1:
		title = "第二題"
	case 2:
		title = "第三題"
	case 3:
		title = "第四題"
	case 4:
		title = "第五題"
	}

	question := Questions[questionID]

	options := make([]linebot.FlexComponent, 0)
	for index, option := range question.Options {
		optionComponent := &linebot.ButtonComponent{
			Type: linebot.FlexComponentTypeButton,
			Action: &linebot.MessageAction{
				Label: option,
				Text:  strconv.Itoa(index + 1),
			},
		}
		options = append(options, optionComponent)
	}

	flexContainer := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: title,
				},
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   question.Question,
					Size:   linebot.FlexTextSizeTypeLg,
					Weight: linebot.FlexTextWeightTypeBold,
					Wrap:   true,
				},
			},
		},
		Body: &linebot.BoxComponent{
			Type:     linebot.FlexComponentTypeBox,
			Layout:   linebot.FlexBoxLayoutTypeVertical,
			Contents: options,
		},
	}

	return linebot.NewFlexMessage(fmt.Sprintf("Question %v", answeredCount+1), flexContainer)
}

func (q *Question) ReasonMessage(answer string) *linebot.FlexMessage {
	var title string

	if q.Answer == answer {
		title = "答對拉 (つ´ω`)つ"
	} else {
		title = "答錯拉 QAQ 可惜了～～～"
	}

	flexContainer := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: title,
					Size: linebot.FlexTextSizeTypeXxl,
				},
			},
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "原因",
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeXl,
				},
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: q.Reason,
					Wrap: true,
				},
			},
		},
	}

	return linebot.NewFlexMessage("result", flexContainer)
}

func (p *Player) FinalMessage() *linebot.FlexMessage {
	flexContainer := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: fmt.Sprintf("你答對了 %v 題!!!", p.Score),
					Size: linebot.FlexTextSizeTypeXxl,
				},
			},
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: "遊戲結束拉 ～～～",
				},
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: "請填一下下面的表單歐 謝謝！",
				},
			},
		},
		Footer: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type: linebot.FlexComponentTypeButton,
					Action: &linebot.URIAction{
						Label: "填寫表單",
						URI:   "http://linecorp.com/",
					},
				},
			},
		},
	}

	return linebot.NewFlexMessage("final score", flexContainer)
}
