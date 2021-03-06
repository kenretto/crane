package filetype

import (
	"bytes"
	"encoding/hex"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Dictionary ...
type Dictionary struct {
	TypeMap sync.Map
	sync.Once
}

var (
	dictionary = new(Dictionary)
	// Unknown unknown file type
	Unknown = "unknown"
)

// FileType according to the content of some files passed in, determine the type of files
func FileType(fSrc []byte) string {
	if fSrc == nil {
		return "unknown"
	}
	typ := dictionary.FileType(fSrc)
	if typ == "" {
		typ = Unknown
	}
	return typ
}

// FileBytes pass in the absolute path of the file and return part of the file
func FileBytes(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}

	var b = make([]byte, 64)
	_, _ = file.Read(b)
	return b
}

func (dictionary *Dictionary) lazyLoad() {
	dictionary.Do(func() {
		dictionary.TypeMap.Store("ffd8ffe", "jpg")                     // JPEG (jpg)
		dictionary.TypeMap.Store("89504e470d0a1a0a0000", "png")        // PNG (png)
		dictionary.TypeMap.Store("474946383961", "gif")                // GIF (gif)
		dictionary.TypeMap.Store("49492a00227105008037", "tif")        // TIFF (tif)
		dictionary.TypeMap.Store("424d228c010000000000", "bmp")        // 16色位图(bmp)
		dictionary.TypeMap.Store("424d8240090000000000", "bmp")        // 24位位图(bmp)
		dictionary.TypeMap.Store("424d8e1b030000000000", "bmp")        // 256色位图(bmp)
		dictionary.TypeMap.Store("41433130313500000000", "dwg")        // CAD (dwg)
		dictionary.TypeMap.Store("3c21444f435459504520", "html")       // HTML (html)   3c68746d6c3e0  3c68746d6c3e0
		dictionary.TypeMap.Store("3c68746d6c3e0", "html")              // HTML (html)   3c68746d6c3e0  3c68746d6c3e0
		dictionary.TypeMap.Store("3c21646f637479706520", "htm")        // HTM (htm)
		dictionary.TypeMap.Store("48544d4c207b0d0a0942", "css")        // css
		dictionary.TypeMap.Store("696b2e71623d696b2e71", "js")         // js
		dictionary.TypeMap.Store("7b5c727466315c616e73", "rtf")        // Rich Text Format (rtf)
		dictionary.TypeMap.Store("38425053000100000000", "psd")        // Photoshop (psd)
		dictionary.TypeMap.Store("46726f6d3a203d3f6762", "eml")        // Email [Outlook Express 6] (eml)
		dictionary.TypeMap.Store("d0cf11e0a1b11ae10000", "doc")        // MS Excel 注意：word、msi 和 excel的文件头一样
		dictionary.TypeMap.Store("d0cf11e0a1b11ae10000", "vsd")        // Visio 绘图
		dictionary.TypeMap.Store("5374616E64617264204A", "mdb")        // MS Access (mdb)
		dictionary.TypeMap.Store("252150532D41646F6265", "ps")         //
		dictionary.TypeMap.Store("255044462d312", "pdf")               // Adobe Acrobat (pdf)
		dictionary.TypeMap.Store("2e524d46000000120001", "rmvb")       // rmvb/rm相同
		dictionary.TypeMap.Store("464c5601050000000900", "flv")        // flv与f4v相同
		dictionary.TypeMap.Store("00000020667479706d70", "mp4")        //
		dictionary.TypeMap.Store("49443303000000002176", "mp3")        //
		dictionary.TypeMap.Store("000001ba210001000180", "mpg")        //
		dictionary.TypeMap.Store("3026b2758e66cf11a6d9", "wmv")        // wmv与asf相同
		dictionary.TypeMap.Store("52494646e27807005741", "wav")        // Wave (wav)
		dictionary.TypeMap.Store("52494646d07d60074156", "avi")        //
		dictionary.TypeMap.Store("4d546864000000060001", "mid")        // MIDI (mid)
		dictionary.TypeMap.Store("504b0304140000000800", "zip")        //
		dictionary.TypeMap.Store("526172211a0700cf9073", "rar")        //
		dictionary.TypeMap.Store("235468697320636f6e66", "ini")        //
		dictionary.TypeMap.Store("504b03040a0000000000", "jar")        //
		dictionary.TypeMap.Store("4d5a9000030000000400", "exe")        // 可执行文件
		dictionary.TypeMap.Store("3c25402070616765206c", "jsp")        // jsp文件
		dictionary.TypeMap.Store("4d616e69666573742d56", "mf")         // MF文件
		dictionary.TypeMap.Store("3c3f786d6c2076657273", "xml")        // xml文件
		dictionary.TypeMap.Store("494e5345525420494e54", "sql")        // xml文件
		dictionary.TypeMap.Store("7061636b616765207765", "java")       // java文件
		dictionary.TypeMap.Store("406563686f206f66660d", "bat")        // bat文件
		dictionary.TypeMap.Store("1f8b0800000000000000", "gz")         // gz文件
		dictionary.TypeMap.Store("6c6f67346a2e726f6f74", "properties") // bat文件
		dictionary.TypeMap.Store("cafebabe0000002e0041", "class")      // bat文件
		dictionary.TypeMap.Store("49545346030000006000", "chm")        // bat文件
		dictionary.TypeMap.Store("04000000010000001300", "mxp")        // bat文件
		dictionary.TypeMap.Store("504b0304140006000800", "docx")       // docx文件
		dictionary.TypeMap.Store("d0cf11e0a1b11ae10000", "wps")        // WPS文字wps、表格et、演示dps都是一样的
		dictionary.TypeMap.Store("6431303a637265617465", "torrent")    //
		dictionary.TypeMap.Store("6D6F6F76", "mov")                    // Quicktime (mov)
		dictionary.TypeMap.Store("FF575043", "wpd")                    // WordPerfect (wpd)
		dictionary.TypeMap.Store("CFAD12FEC5FD746F", "dbx")            // Outlook Express (dbx)
		dictionary.TypeMap.Store("2142444E", "pst")                    // Outlook (pst)
		dictionary.TypeMap.Store("AC9EBD8F", "qdf")                    // Quicken (qdf)
		dictionary.TypeMap.Store("E3828596", "pwl")                    // Windows Password (pwl)
		dictionary.TypeMap.Store("2E7261FD", "ram")                    // Real Audio (ram)
	})
}

// 获取前面结果字节的二进制
func (dictionary *Dictionary) bytesToHexString(src []byte) string {
	res := bytes.Buffer{}
	if src == nil || len(src) <= 0 {
		return ""
	}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10))
		}
		res.WriteString(hv)
	}
	return res.String()
}

// FileType 用文件前面几个字节来判断
// fSrc: 文件字节流（就用前面几个字节）
func (dictionary *Dictionary) FileType(fSrc []byte) string {
	dictionary.lazyLoad()
	var fileType string
	fileCode := dictionary.bytesToHexString(fSrc)
	dictionary.TypeMap.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		if strings.HasPrefix(fileCode, strings.ToLower(k)) ||
			strings.HasPrefix(k, strings.ToLower(fileCode)) {
			fileType = v
			return false
		}
		return true
	})
	return fileType
}
