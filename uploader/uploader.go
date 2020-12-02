package uploader

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kenretto/crane/filetype"
	"mime/multipart"
)

// SaveHandler 自定义文件上传之后的保存操作
type SaveHandler interface {
	// 保存文件并返回文件最终路径
	Save(file *multipart.FileHeader, fileName string) (string, error)
}

// DefaultSaveHandler 默认文件保存器
type DefaultSaveHandler struct {
	prefix  string
	dst     string
	context *gin.Context
}

// SetDst set save file dir
func (defaultSaveHandler *DefaultSaveHandler) SetDst(dst string) *DefaultSaveHandler {
	defaultSaveHandler.dst = dst
	return defaultSaveHandler
}

// SetPrefix set save file prefix
func (defaultSaveHandler *DefaultSaveHandler) SetPrefix(prefix string) *DefaultSaveHandler {
	defaultSaveHandler.prefix = prefix
	return defaultSaveHandler
}

// Save save
func (defaultSaveHandler *DefaultSaveHandler) Save(file *multipart.FileHeader, fileName string) (string, error) {
	filePath := defaultSaveHandler.dst + defaultSaveHandler.prefix + fileName
	err := defaultSaveHandler.context.SaveUploadedFile(file, filePath)
	if err != nil {
		return "", err
	}

	return filePath, err
}

// Uploader file upload function
type Uploader struct {
	FormKey      string
	SaveHandler  SaveHandler
	AllowedTypes []string
	NameFn       func(index int, file *multipart.FileHeader) string
	Ctx          *gin.Context
}

// TypeValid valid file type
func (u *Uploader) TypeValid(file *multipart.FileHeader) error {
	for _, typ := range u.AllowedTypes {
		f, err := file.Open()
		if err != nil {
			return err
		}

		var b = make([]byte, 64)
		_, _ = f.Read(b)
		if typ == filetype.FileType(b) {
			return nil
		}
	}
	return errors.New("file type not allowed")
}

// Save save file
func (u *Uploader) Save(index int, file *multipart.FileHeader) (filename string, err error) {
	err = u.TypeValid(file)
	if err != nil {
		return
	}
	filename, err = u.SaveHandler.Save(file, u.NameFn(index, file))
	return
}

// Files get all uploaded files
func (u *Uploader) Files() []*multipart.FileHeader {
	form, err := u.Ctx.MultipartForm()
	if err != nil {
		return nil
	}
	files := form.File[u.FormKey]
	return files
}

// Each perform traversal processing operations on uploaded files
func (u *Uploader) Each(fn func(index int, file *multipart.FileHeader) error) error {
	for i, header := range u.Files() {
		err := fn(i, header)
		if err != nil {
			return err
		}
	}
	return nil
}

// SaveAll 文件上传
//  key 上传文件的表单name, 如果是多文件需要加上中括号[]
//  dst 存放路径 注意:无论这里传什么路径, 最后边都会追加 filename.xxx
func (u *Uploader) SaveAll() (files []string, err error) {
	err = u.Each(func(index int, file *multipart.FileHeader) error {
		filename, err := u.Save(index, file)
		if err != nil {
			return err
		}
		files = append(files, filename)
		return nil
	})
	return
}
