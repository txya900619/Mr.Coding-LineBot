package mrcoding

import (
	"Mr.Coding-LineBot/spreadsheets"
	"github.com/line/line-bot-sdk-go/linebot"
)

func getCompleteFormFlexContainer(chatroomID string) linebot.FlexContainer {
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Size: linebot.FlexBubbleSizeTypeKilo,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "Mr.Coding 聊天室",
					Size:   linebot.FlexTextSizeTypeXl,
					Weight: linebot.FlexTextWeightTypeBold,
				},
			},
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: "恭喜你填完表單！",
				},
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: "立即點擊下方按鈕進入聊天室～",
				},
			},
		},
		Footer: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type:   linebot.FlexComponentTypeButton,
					Height: linebot.FlexButtonHeightTypeSm,
					Action: &linebot.URIAction{
						Label: "點我進入",
						// Liff url
						URI: "https://liff.line.me/1654852713-gR9j0RyE/" + chatroomID,
					},
				},
			},
		},
		Styles: &linebot.BubbleStyle{
			Footer: &linebot.BlockStyle{Separator: true},
		},
	}
}

func getQuestionFlexContainer(questionID spreadsheets.ColumnID) linebot.FlexContainer {
	text := ""
	instructions := ""
	footerPassBtn := false
	switch questionID {
	case spreadsheets.QuestionEmail:
		text = "輸入您的電子郵件地址"
		instructions = "限定輸入一行"
	case spreadsheets.QuestionName:
		text = "輸入您的姓名"
		instructions = "限定輸入一行"
	case spreadsheets.QuestionStudentNo:
		text = "輸入您的學號"
		instructions = "限定輸入一行，例如：107000001"
	case spreadsheets.QuestionDepartment:
		text = "輸入您的系級"
		instructions = "限定輸入一行，例如：電資一"
	case spreadsheets.QuestionProgramming:
		text = "輸入您的程式問題"
		instructions = "允許多行輸入"
	case spreadsheets.QuestionUploadFile:
		text = "上傳程式截圖或程式碼網址"
		instructions = "直接上傳圖片檔，僅限圖片檔"
		footerPassBtn = true
	case spreadsheets.QuestionNote:
		text = "輸入其他您想說的"
		instructions = "允許多行輸入"
		footerPassBtn = true

	}

	container := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Size: linebot.FlexBubbleSizeTypeKilo,
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   text,
					Size:   linebot.FlexTextSizeTypeXl,
					Weight: linebot.FlexTextWeightTypeBold,
					Wrap:   true,
				},
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   instructions,
					Size:   linebot.FlexTextSizeTypeSm,
					Color:  "#808080",
					Margin: linebot.FlexComponentMarginTypeMd,
				},
			},
		},
	}

	if footerPassBtn {
		container.Footer = &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type:   linebot.FlexComponentTypeButton,
					Height: linebot.FlexButtonHeightTypeSm,
					Action: &linebot.PostbackAction{
						Label:       "略過此題",
						Data:        "pass",
						DisplayText: "略過",
					},
				},
			},
		}
		container.Styles = &linebot.BubbleStyle{
			Footer: &linebot.BlockStyle{Separator: true},
		}
	}
	return container
}

func getEmailErrorFlexContainer() linebot.FlexContainer {
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: "請輸入合法格式的 email",
				},
			},
		},
	}
}

func HelpMessage() *linebot.FlexMessage {
	flexContainer := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type: linebot.FlexComponentTypeButton,
					Action: &linebot.MessageAction{
						Label: "點我填寫 Mr.Coding 表單",
						Text:  "Mr.Coding 表單",
					},
				},
				&linebot.ButtonComponent{
					Type: linebot.FlexComponentTypeButton,
					Action: &linebot.MessageAction{
						Label: "點我遊玩社團博覽會有獎徵答",
						Text:  "社團博覽會有獎徵答",
					},
				},
			},
		},
	}

	return linebot.NewFlexMessage("help", flexContainer)
}

func GetTypeErrorFlexContainer() linebot.FlexContainer {
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: "請不要輸入文字以外的訊息",
				},
			},
		},
	}
}

func GetWorkingFlexContainer() linebot.FlexContainer {
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: "目前正在維護中，預計 10 月重啟服務",
				},
			},
		},
	}
}
