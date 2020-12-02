package category

import (
	"fmt"
	"github.com/kenretto/crane/database"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

type (
	// Node classification information
	Node struct {
		Alias       string `json:"alias"`       // It may be the data needed for page link splicing, such as items?category={alias} or items/{alias}
		Name        string `json:"name"`        // category display name
		Pic         string `json:"pic"`         // icon
		Badge       string `json:"badge"`       // badges, possible corner markers
		Description string `json:"description"` // possible introduction
		Next        Nodes  `json:"next"`        // subcategory
	}
	// Nodes specify a level of classification information
	Nodes []*Node

	// Tree all classified raw data (sorted)
	Tree struct {
		items    reflect.Value
		category *Category
	}

	// Category classification
	Category struct {
		aliasField, nameField, picField, badgeField, levelField, descriptionField, PidField, IDField string
		tableName                                                                                    func() string
		tableTyp                                                                                     reflect.Type
	}

	// Table orm table handler
	Table struct {
		category *Category
	}

	// Option setup
	Option func(category *Category)
)

var (
	defaultSetting = []Option{
		SetDescriptionField("Description"), SetAliasField("Alias"), SetBadgeField("Badge"),
		SetLevelField("Level"), SetNameField("Name"), SetPicField("Pic"), SetPidField("Pid"), SetIDField("ID"),
	}
)

// SetDescriptionField set the field name of the detail profile
func SetDescriptionField(field string) Option {
	return func(category *Category) {
		if category.descriptionField == "" {
			category.descriptionField = field
		}
	}
}

// SetPidField set parent identity field name
func SetPidField(field string) Option {
	return func(category *Category) {
		if category.PidField == "" {
			category.PidField = field
		}
	}
}

// SetIDField set id field name
func SetIDField(field string) Option {
	return func(category *Category) {
		if category.IDField == "" {
			category.IDField = field
		}
	}
}

// SetLevelField set hierarchical field name
func SetLevelField(field string) Option {
	return func(category *Category) {
		if category.levelField == "" {
			category.levelField = field
		}
	}
}

// SetBadgeField set corner badge field name
func SetBadgeField(field string) Option {
	return func(category *Category) {
		if category.badgeField == "" {
			category.badgeField = field
		}
	}
}

// SetPicField set icon field name
func SetPicField(field string) Option {
	return func(category *Category) {
		if category.picField == "" {
			category.picField = field
		}
	}
}

// SetNameField set display name field name
func SetNameField(field string) Option {
	return func(category *Category) {
		if category.nameField == "" {
			category.nameField = field
		}
	}
}

// SetAliasField set alias field field name
func SetAliasField(field string) Option {
	return func(category *Category) {
		if category.aliasField == "" {
			category.aliasField = field
		}
	}
}

func newItems(typ reflect.Type) reflect.Value {
	slice := reflect.MakeSlice(reflect.SliceOf(typ), 0, 0)
	items := reflect.New(slice.Type())
	items.Elem().Set(slice)
	return items
}

// New classification
func New(options ...Option) *Category {
	var category = new(Category)
	for _, option := range options {
		option(category)
	}

	for _, option := range defaultSetting {
		option(category)
	}

	return category
}

// Table get the table structure do not pass a pointer
func (category *Category) Table(table database.Table) *Table {
	category.tableName = table.TableName
	category.tableTyp = reflect.TypeOf(table)
	return &Table{category: category}
}

// WithDB use gorm
func (table *Table) WithDB(db *gorm.DB) *Tree {
	items := newItems(table.category.tableTyp)
	field, _ := table.category.tableTyp.FieldByName(table.category.levelField)
	column := func() string {
		tags := strings.Split(field.Tag.Get("gorm"), ",")
		for _, tag := range tags {
			tagInfo := strings.Split(tag, ":")
			if len(tagInfo) == 2 && strings.ToLower(tagInfo[0]) == "column" {
				return tagInfo[1]
			}
		}
		return table.category.levelField
	}()

	db.Table(table.category.tableName()).Order(fmt.Sprintf("%s asc", column)).Find(items.Interface())
	return &Tree{items: items.Elem(), category: table.category}
}

// Categories get tree
func (tree *Tree) Categories() Nodes {
	dataLen := tree.items.Len()
	var nodes, tmp = make(Nodes, 0), make(map[string]*Node)

	for i := 0; i < dataLen; i++ {
		item := tree.items.Index(i)
		pid, id := item.FieldByName(tree.category.PidField).String(), item.FieldByName(tree.category.IDField).String()
		tmp[id] = newNode(item, tree.category)

		if pid == "" {
			// If pid is empty, it is the top-level classification by default
			nodes = append(nodes, tmp[id])
		} else {
			if _, ok := tmp[pid]; ok {
				if tmp[pid].Next == nil {
					tmp[pid].Next = make(Nodes, 0)
				}
				// 如果在pos中存在父id为pid的数据,直接将这次的数据追加到pos中该组数据中的children中
				tmp[pid].Next = append(tmp[pid].Next, tmp[id])
			} else {
				// 否则默认无父分类，直接写入tree中
				nodes = append(nodes, tmp[id])
			}
		}
	}

	return nodes
}

func mustHasValue(value reflect.Value) string {
	v := value.Interface()
	switch v := v.(type) {
	case string:
		return v
	}
	return ""
}

func newNode(item reflect.Value, category *Category) *Node {
	return &Node{
		Alias:       mustHasValue(item.FieldByName(category.aliasField)),
		Name:        mustHasValue(item.FieldByName(category.nameField)),
		Pic:         mustHasValue(item.FieldByName(category.picField)),
		Badge:       mustHasValue(item.FieldByName(category.badgeField)),
		Description: mustHasValue(item.FieldByName(category.descriptionField)),
	}
}
