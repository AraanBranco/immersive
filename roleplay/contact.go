package roleplay

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/dgraph-io/badger"
	badgerhold "github.com/timshannon/badgerhold"
)

type Contact struct {
	ID   string `badgerholdIndex:"IdxContact"`
	Name string
	Org  string
	Dw   string
}

func ContactCommand(ctx Context) {
	if len(ctx.Args) == 0 {
		msg := fmt.Sprintf("Comandos disponíveis para contato: \n")
		msg += fmt.Sprintf("`!contato create <id> <nome> <org>` \n")
		msg += fmt.Sprintf("`!contato update <id> <nome> <org>` \n")
		msg += fmt.Sprintf("`!contato search <id>`")
		ctx.Reply(msg)
		return
	}

	// Create name for channel
	if ctx.Args[0] == "create" {
		// Valid if pass ID and Name
		if len(ctx.Args) < 3 {
			ctx.Reply("Para criar um novo contato siga o exemplo: `!contato create <id>* <nome>* <org>` **(*)Campos obrigatórios**")
			return
		}

		str := createContact(ctx)
		ctx.Reply(str)
		return
	}

	if ctx.Args[0] == "update" {
		if len(ctx.Args) == 1 {
			ctx.Reply("Para atualizar um contato siga o exemplo: `!contato update <id>* <nome>* <org>` **(*)Campos obrigatórios**")
			return
		}

		str := updateContact(ctx)
		ctx.Reply(str)
		return
	}

	if ctx.Args[0] == "search" {
		if len(ctx.Args) == 1 {
			ctx.Reply("Para buscar um contato siga o exemplo: `!contato search <id>` **Passe o passaporte**")
			return
		}

		str := searchContact(ctx)
		ctx.Reply(str)
		return
	}
}

func createContact(ctx Context) string {
	c := &Contact{
		ID:   ctx.Args[1],
		Name: strings.Title(ctx.Args[2]),
		Org:  strings.Title(ctx.Args[3]),
		Dw:   createDw(strings.ToLower(ctx.Args[2])),
	}

	err := createContactDB(ctx, *c)
	if err != nil {
		fmt.Println("Erro ao cadastrar contato: ", err)
		return "Ocorreu um erro ao cadastrar o contato!"
	}

	// Mount layout for contact
	contact := fmt.Sprintf("Passaporte: %v \n", c.ID)
	contact += fmt.Sprintf("Nome: %s \n", c.Name)
	contact += fmt.Sprintf("Org: %s \n", c.Org)
	contact += fmt.Sprintf("DeepWeb: %s", c.Dw)

	buffer := bytes.NewBufferString("```")
	buffer.WriteString(contact)
	buffer.WriteString("```")

	return buffer.String()
}

func createContactDB(ctx Context, c Contact) error {
	store := ctx.DB
	err := store.Badger().Update(func(tx *badger.Txn) error {
		err := store.TxInsert(tx, c.ID, c)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func updateContact(ctx Context) string {
	search, err := searchContactDB(ctx, ctx.Args[1])
	if len(search) == 0 || err != nil {
		return "Nenhum contato encontrado!"
	}

	name := search[0].Name
	org := search[0].Org

	if len(ctx.Args) >= 3 && ctx.Args[2] != "" {
		name = ctx.Args[2]
	}

	if len(ctx.Args) >= 4 && ctx.Args[3] != "" {
		org = ctx.Args[3]
	}

	c := &Contact{
		ID:   ctx.Args[1],
		Name: strings.Title(name),
		Org:  strings.Title(org),
		Dw:   search[0].Dw,
	}

	err = updateContactDB(ctx, *c)
	if err != nil {
		fmt.Println("Erro ao cadastrar contato: ", err)
		return "Ocorreu um erro ao cadastrar o contato!"
	}

	// Mount layout for contact
	contact := fmt.Sprintf("Passaporte: %v \n", c.ID)
	contact += fmt.Sprintf("Nome: %s \n", c.Name)
	contact += fmt.Sprintf("Org: %s \n", c.Org)
	contact += fmt.Sprintf("DeepWeb: %s", c.Dw)

	buffer := bytes.NewBufferString("```")
	buffer.WriteString(contact)
	buffer.WriteString("```")

	return buffer.String()
}

func updateContactDB(ctx Context, c Contact) error {
	store := ctx.DB
	err := store.Update(c.ID, c)

	return err
}

func searchContact(ctx Context) string {
	param := ctx.Args[1]
	contact := "Nenhum contato localizado!"
	results, err := searchContactDB(ctx, param)
	if err != nil {
		fmt.Println("Erro ao buscar contato: ", err)
		return "Ocorreu um erro ao buscar contato!"
	}

	if len(results) > 0 {
		// Mount layout for contact
		contact = "Contato localizado!\n"
		contact += "```"
		contact += fmt.Sprintf("Passaporte: %v \n", results[0].ID)
		contact += fmt.Sprintf("Nome: %s \n", results[0].Name)
		contact += fmt.Sprintf("Org: %s \n", results[0].Org)
		contact += fmt.Sprintf("DeepWeb: %s", results[0].Dw)
		contact += "```"
	}

	buffer := bytes.NewBufferString(contact)
	return buffer.String()
}

func searchContactDB(ctx Context, param string) ([]Contact, error) {
	store := ctx.DB
	var contact []Contact
	err := store.Find(&contact, badgerhold.Where("ID").Eq(param))
	return contact, err
}

func createDw(name string) string {
	return fmt.Sprintf("%s-%v", name, randomInt())
}

func randomInt() int {
	rand.Seed(time.Now().UnixNano())
	min := 0
	max := 9999

	return rand.Intn(max-min) + min
}
