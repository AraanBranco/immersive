package roleplay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"sort"
	"strings"
)

type Outfit struct {
	Mecanico []DefaultStruct `json:"mecanico,omitempty"`
	Motoclub []DefaultStruct `json:"motoclub,omitempty"`
	Ballas   []DefaultStruct `json:"ballas,omitempty"`
	Corleone []DefaultStruct `json:"corleone,omitempty"`
}

type DefaultStruct struct {
	Acessorios string `json:"acessorios,omitempty"`
	Chapeu     string `json:"chapeu,omitempty"`
	Oculos     string `json:"oculos,omitempty"`
	Mochila    string `json:"mochila,omitempty"`
	Jaqueta    string `json:"jaqueta,omitempty"`
	Calca      string `json:"calca,omitempty"`
	Sapatos    string `json:"sapatos,omitempty"`
	Maos       string `json:"maos,omitempty"`
	Imagem     string `json:"imagem,omitempty"`
}

func OutfitCommand(ctx Context) {
	if len(ctx.Args) == 0 {
		ctx.Reply("Para saber mais sobre os comandos digite: `!outfit <cidade> <org>`")
		return
	} else {
		// Get struct for city
		outfits := getOutfit(ctx.Args[0])
		allowCommands := allCommands(*outfits)
		sort.Strings(allowCommands)

		// Check if exist ORG
		if !contains(allowCommands, ctx.Args[1]) {
			msg := fmt.Sprintf("Outfit não existe, tente um desses: `%s`", strings.Join(allowCommands[:], ", "))
			ctx.Reply(msg)
			return
		}

		selectOutfit := strings.Title(ctx.Args[1])
		rv := reflect.ValueOf(&outfits).Elem().Elem().FieldByName(selectOutfit).Addr()
		outfit := rv.Interface().(*[]DefaultStruct)
		// Valid if exist outfir for ORG
		if len(*outfit) == 0 {
			ctx.Reply("Nenhum outfit disponível para essa ORG")
			return
		}

		// Mount messages with outfits
		ctx.Reply(fmt.Sprintf("Outfits para %s: \n", selectOutfit))
		for _, item := range *outfit {
			buff := bytes.NewBufferString("```")
			allStruct := reflect.ValueOf(item)
			img := ""
			for i := 0; i < allStruct.NumField(); i++ {
				field := allStruct.Type().Field(i)
				value := allStruct.Field(i)
				msg := ""
				if value.Interface() != "" && field.Name != "Imagem" {
					msg = fmt.Sprintf("%s: %s\n", field.Name, value.Interface())
				}

				if field.Name == "Imagem" {
					img = fmt.Sprintf("%s \n", value.Interface())
				}
				buff.WriteString(msg)
			}
			buff.WriteString("```")
			buff.WriteString(img)
			buff.WriteString("------------------------")

			str2 := buff.String()
			ctx.Reply(str2)
		}
	}
}

func allCommands(inter interface{}) []string {
	all := reflect.TypeOf(inter)
	var result []string
	for i := 0; i < all.NumField(); i++ {
		name := strings.ToLower(all.Field(i).Name)
		result = append(result, name)
	}
	return result
}

func getOutfit(city string) *Outfit {
	pathFile := fmt.Sprintf("outfits/%s.json", city)
	body, err := ioutil.ReadFile(pathFile)
	if err != nil {
		fmt.Println("Erro ao carregar outfit da cidade: ", err)
		return nil
	}
	var outfit Outfit
	json.Unmarshal(body, &outfit)
	return &outfit
}

func contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}
