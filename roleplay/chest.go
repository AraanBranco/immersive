package roleplay

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dgraph-io/badger"
	badgerhold "github.com/timshannon/badgerhold"
)

type Chest struct {
	Name    string `badgerholdIndex:"Item"`
	Qtd     int
	Created time.Time
	Updated time.Time
}

var today = time.Now().UTC()

func ChestCommand(ctx Context) {
	if len(ctx.Args) == 0 {
		msg := fmt.Sprintf("Comandos disponíveis para chest: \n")
		msg += fmt.Sprintf("`!chest add <item> <qtd>` \n")
		msg += fmt.Sprintf("`!chest rm <item> <qtd>` \n")
		msg += fmt.Sprintf("`!chest list`")
		ctx.Reply(msg)
		return
	}

	// Add item in Chest
	if ctx.Args[0] == "add" {
		// Valid if pass ID and Name
		if len(ctx.Args) < 2 {
			ctx.Reply("Para criar um novo contato siga o exemplo: `!chest add <item>* <qtd>*` **(*)Campos obrigatórios**")
			return
		}

		str := addItem(ctx)
		ctx.Reply(str)
		return
	}

	// Remove item from chest
	if ctx.Args[0] == "rm" {
		if len(ctx.Args) == 1 {
			ctx.Reply("Para retirar um item siga o exemplo: `!chest rm <item>* <qtd>*` **(*)Campos obrigatórios**")
			return
		}

		str := rmItem(ctx)
		ctx.Reply(str)
		return
	}

	// List itens in Chest
	if ctx.Args[0] == "list" {
		str := listItem(ctx)
		ctx.Reply(str)
		return
	}
}

func addItem(ctx Context) string {
	nameItem := ctx.Args[1]
	comm := "modificado"
	var c Chest

	search, err := getItemDB(ctx, nameItem)
	if len(search) == 0 || err != nil {
		qtd, _ := strconv.Atoi(ctx.Args[2])
		c = Chest{
			Name:    strings.ToLower(nameItem),
			Qtd:     qtd,
			Created: today,
		}
		comm = "adicionado"

		err = addItemDB(ctx, c)
		if err != nil {
			fmt.Println("Erro ao cadastrar contato: ", err)
			return "Ocorreu um erro ao adicionar o item ao chest!"
		}
	} else {
		qtd, _ := strconv.Atoi(ctx.Args[2])
		qtd = qtd + search[0].Qtd
		c = Chest{
			Name:    search[0].Name,
			Qtd:     qtd,
			Updated: today,
		}

		err = updateItemDB(ctx, c)
		if err != nil {
			fmt.Println("Erro ao cadastrar contato: ", err)
			return "Ocorreu um erro ao atualizar o item ao chest!"
		}
	}

	// Mount layout for Chest
	msg := fmt.Sprintf("Item %v: %v - Quantidade: %v \n", comm, c.Name, c.Qtd)

	buffer := bytes.NewBufferString("```")
	buffer.WriteString(msg)
	buffer.WriteString("```")

	return buffer.String()
}

func addItemDB(ctx Context, c Chest) error {
	store := ctx.DB
	err := store.Badger().Update(func(tx *badger.Txn) error {
		err := store.TxInsert(tx, c.Name, c)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func rmItem(ctx Context) string {
	var msg string
	search, err := getItemDB(ctx, ctx.Args[1])
	if len(search) == 0 || err != nil {
		return "Nenhum item encontrado no chest!"
	}

	name := search[0].Name
	qtd := search[0].Qtd

	// Remove item from Chest
	if len(ctx.Args) == 2 {
		err = rmItemDB(ctx, Chest{Name: name})
		if err != nil {
			fmt.Println("Erro ao remover item: ", err)
			return "Ocorreu um erro ao remover o item!"
		}
		msg = fmt.Sprintf("``` %v removido do Chest ```", name)

		// Send Log to owner discord
		LogChestToOwner(ctx, msg)
		return msg
	}

	if len(ctx.Args) == 3 {
		qtdAux, _ := strconv.Atoi(ctx.Args[2])
		aux := qtd - qtdAux

		if aux == 0 {
			err = rmItemDB(ctx, Chest{Name: name})
			if err != nil {
				fmt.Println("Erro ao remover item: ", err)
				return "Ocorreu um erro ao remover o item!"
			}
			msg = fmt.Sprintf("``` %v removido do Chest ```", name)

			// Send Log to owner discord
			LogChestToOwner(ctx, msg)
			return fmt.Sprintf(msg)
		}

		c := &Chest{
			Name: strings.ToLower(name),
			Qtd:  aux,
		}

		err = updateItemDB(ctx, *c)
		if err != nil {
			fmt.Println("Erro ao cadastrar contato: ", err)
			return "Ocorreu um erro ao atualizar o item ao chest!"
		}

		// Mount layout for contact
		msg = fmt.Sprintf("Item modificado: %v - Quantidade Total: %v \n", c.Name, c.Qtd)

		buffer := bytes.NewBufferString("```")
		buffer.WriteString(msg)
		buffer.WriteString("```")

		return buffer.String()
	}

	return "Tente novamente!"
}

func rmItemDB(ctx Context, c Chest) error {
	store := ctx.DB
	err := store.DeleteMatching(c, badgerhold.Where("Name").Eq(c.Name))

	return err
}

func listItem(ctx Context) string {
	msg := "O chest está vazio!"
	results, err := listItemDB(ctx)
	if err != nil {
		fmt.Println("Erro ao abrir o chest: ", err)
		return "Ocorreu um erro ao abrir o chest!"
	}

	if len(results) > 0 {
		// Mount layout for chest
		msg = "Itens no chest!\n"
		msg += "```"
		msg += "ITEM | QUANTIDADE \n"
		msg += fmt.Sprintf("--------------------------------------------- \n")
		for _, item := range results {
			msg += fmt.Sprintf("%v | %v \n", strings.Title(item.Name), item.Qtd)
			msg += fmt.Sprintf("--------------------------------------------- \n")
		}
		msg += "```"
	}

	buffer := bytes.NewBufferString(msg)
	return buffer.String()
}

func listItemDB(ctx Context) ([]Chest, error) {
	store := ctx.DB
	var chest []Chest
	err := store.Find(&chest, badgerhold.Where("Created").Ne(today).SortBy("Name"))
	return chest, err
}

func updateItemDB(ctx Context, c Chest) error {
	store := ctx.DB
	err := store.Update(c.Name, c)

	return err
}

func getItemDB(ctx Context, param string) ([]Chest, error) {
	store := ctx.DB
	var chest []Chest
	err := store.Find(&chest, badgerhold.Where("Name").Eq(param))
	return chest, err
}
