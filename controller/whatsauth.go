package controller

import (
	"gocroot/config"
	"gocroot/helper"
	"gocroot/model"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func WhatsAuthReceiver(c *fiber.Ctx) error {
	var h model.Header
	err := c.ReqHeaderParser(&h)
	if err != nil {
		return err
	}
	var resp model.Response
	if h.Secret == config.WebhookSecret {
		var msg model.IteungMessage
		err = c.BodyParser(&msg)
		if err != nil {
			return err
		}
		if IsLoginRequest(msg, config.WAKeyword) { 
			resp = HandlerQRLogin(msg, config.WAKeyword)
		} else { 
			resp = HandlerIncomingMessage(msg)
		}
	} else {
		resp.Response = "Secret Salah"
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func RefreshWAToken(c *fiber.Ctx) error {
	dt := &model.WebHook{
		URL:    config.WebhookURL,
		Secret: config.WebhookSecret,
	}
	resp, err := helper.PostStructWithToken[model.User]("Token", WAAPIToken(config.WAPhoneNumber), dt, config.WAAPIGetToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "response": resp})
	}
	profile := &model.Profile{
		Phonenumber: resp.PhoneNumber,
		Token:       resp.Token,
	}
	res, err := helper.ReplaceOneDoc(config.Mongoconn, "profile", bson.M{"phonenumber": resp.PhoneNumber}, profile)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "result": res})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"result": res})
}

func IsLoginRequest(msg model.IteungMessage, keyword string) bool {
	return strings.Contains(msg.Message, keyword) && msg.From_link
}

func GetUUID(msg model.IteungMessage, keyword string) string {
	return strings.Replace(msg.Message, keyword, "", 1)
}

func HandlerQRLogin(msg model.IteungMessage, WAKeyword string) (resp model.Response) {
	dt := &model.WhatsauthRequest{
		Uuid:        GetUUID(msg, WAKeyword),
		Phonenumber: msg.Phone_number,
		Delay:       msg.From_link_delay,
	}
	resp, _ = helper.PostStructWithToken[model.Response]("Token", WAAPIToken(config.WAPhoneNumber), dt, config.WAAPIQRLogin)
	return
}

func HandlerIncomingMessage(msg model.IteungMessage) (resp model.Response) {
	dt := &model.TextMessage{
		To:       msg.Chat_number,
		IsGroup:  false,
		Messages: GetRandomReplyFromMongo(msg),
	}
	if msg.Chat_server == "g.us" {
		dt.IsGroup = true
	}
	if (msg.Phone_number != "628112000279") && (msg.Phone_number != "6283131895000") { //ignore pesan datang dari iteung
		resp, _ = helper.PostStructWithToken[model.Response]("Token", WAAPIToken(config.WAPhoneNumber), dt, config.WAAPIMessage)
	}
	return
}

func GetRandomReplyFromMongo(msg model.IteungMessage) string {
	rply, _ := helper.GetRandomDoc[model.Reply](config.Mongoconn, "reply", 1)
	replymsg := strings.ReplaceAll(rply[0].Message, "#BOTNAME#", msg.Alias_name)
	replymsg = strings.ReplaceAll(replymsg, "\\n", "\n")
	return replymsg
}

func WAAPIToken(phonenumber string) string {
	filter := bson.M{"phonenumber": phonenumber}
	apitoken, _ := helper.GetOneDoc[model.Profile](config.Mongoconn, "profile", filter)
	return apitoken.Token
}
