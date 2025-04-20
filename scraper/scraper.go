package scraper

import (
	"context"
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/ignisVeneficus/library/db"
	"github.com/ignisVeneficus/library/db/dao"
	"github.com/ignisVeneficus/library/db/dbo"
	"github.com/ignisVeneficus/library/utils"

	"github.com/ignisVeneficus/ebook/eBookData"
	"github.com/ignisVeneficus/ebook/epub"
	"github.com/ignisVeneficus/ebook/mobipocket"
)

var (
	ErrDatabase  = errors.New("database error")
	ErrFileError = errors.New("book file error")
)

func readMobiMetadata(filename string) (eBookData.Metadata, []byte, error) {
	var (
		f   *os.File
		err error
	)
	if f, err = os.Open(filename); err != nil {
		log.Logger.Error().Str("file", filename).Err(err).Msg("Reading Error")
		return nil, nil, err
	}
	defer f.Close()
	var mobi *mobipocket.Mobipocket
	if mobi, err = mobipocket.ReadMobi(f); err != nil {
		log.Logger.Error().Str("file", filename).Err(err).Msg("Reading Error")
		return nil, nil, err
	}
	return mobi.Metadata(), mobi.Cover(), nil
}

func readEpubMetadata(fileName string) (eBookData.Metadata, []byte, error) {
	var (
		f   *os.File
		err error
	)
	if f, err = os.Open(fileName); err != nil {
		log.Logger.Error().Str("file", fileName).Err(err).Msg("Reading Error")
		return nil, nil, err
	}
	defer f.Close()
	var epublic *epub.Epub
	if epublic, err = epub.ReadEpub(f, ""); err != nil {
		log.Logger.Error().Str("file", fileName).Err(err).Msg("Reading Error")
		return nil, nil, err
	}
	return epublic.Metadata(), epublic.Cover(), nil
}
func getHash(fileName string) (string, error) {
	var (
		f   *os.File
		err error
	)
	if f, err = os.Open(fileName); err != nil {
		log.Logger.Error().Str("file", fileName).Err(err).Msg("Reading Error")
		return "", err
	}
	hash := md5.New()
	_, err = io.Copy(hash, f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil

}
func getRelativePath(basePath string, fileName string) (string, error) {
	return filepath.Rel(basePath, fileName)
}

func createBook(basePath string, coverBasePath string, fileName string, hash string) (dbo.Book, error) {
	var (
		err         error
		metadata    eBookData.Metadata
		coverData   []byte
		coverFormat string
		bookFormat  string
		color       string
	)
	ext := filepath.Ext(fileName)
	if strings.EqualFold(ext, ".prc") || strings.EqualFold(ext, ".mobi") {
		ext = ".mobi"
		metadata, coverData, err = readMobiMetadata(fileName)
		bookFormat = "mobi"
	} else {
		metadata, coverData, err = readEpubMetadata(fileName)
		bookFormat = "epub"
	}
	if err != nil {
		return dbo.Book{}, err
	}
	hasCover := 0
	if len(coverData) > 0 {
		hasCover = 1
	}
	if hasCover > 0 {
		coverFormat, color, err = ParseCover(coverData)
		if err != nil {
			hasCover = 0
		}
	}
	relativePath, err := getRelativePath(basePath, fileName)
	if err != nil {
		return dbo.Book{}, err
	}
	relativePath = filepath.ToSlash(relativePath)
	isbn := sql.NullString{Valid: false}
	if metadata.ISBN() != "" {
		isbn = sql.NullString{String: metadata.ISBN(), Valid: true}
	}

	book := dbo.Book{
		Title:    sql.NullString{Valid: true, String: metadata.Title()},
		Format:   bookFormat,
		Hash:     hash,
		Hascover: int32(hasCover),
		File:     relativePath,
		Isbn:     isbn,
	}
	authors := make([]dbo.Author, len(metadata.Author()))
	for i, authName := range metadata.Author() {
		authors[i] = dbo.Author{Name: authName}
	}
	book.Authors = authors
	if hasCover > 0 {
		book.CoverColor = sql.NullString{String: color, Valid: true}
		book.CoverType = sql.NullString{String: coverFormat, Valid: true}
		coverFile := filepath.Join(coverBasePath, relativePath+"."+coverFormat)
		coverDir := filepath.Dir(coverFile)
		log.Logger.Trace().Str("file", fileName).Str("Hash", hash).Str("cover dir", coverDir).Msg("Extract Cover")

		err = os.MkdirAll(coverDir, 0777)
		if err != nil {
			log.Logger.Error().Str("file", fileName).Str("Hash", hash).Err(err).Msg("Directory error")
			return book, nil
		}
		err = os.WriteFile(coverFile, coverData, 0666)
		if err != nil {
			log.Logger.Error().Str("file", fileName).Str("Hash", hash).Err(err).Msg("File error")
			return book, nil
		}
	}

	return book, nil
}

func readAFile(observer Observer, basePath string, coverPath string, fileName string, forceUpdate bool) (int64, error) {
	var (
		err      error
		fileBook dbo.Book
	)
	log.Logger.Debug().Str("file", fileName).Msg("Start File reading")
	defer log.Logger.Debug().Str("file", fileName).Msg("End File reading")
	hash, err := getHash(fileName)
	if err != nil {
		log.Logger.Error().Str("file", fileName).Err(err).Msg("File Hash error")
		observer.Error(fileName, "Parse Error")
		return 0, fmt.Errorf("hash error: %w", ErrFileError)
	}
	database := db.GetDatabase()
	ctx := context.Background()
	if forceUpdate {
		fileBook, err = createBook(basePath, coverPath, fileName, hash)
		if err != nil {
			log.Logger.Error().Str("file", fileName).Str("Hash", hash).Err(err).Msg("File error")
			observer.Error(fileName, "Parse Error")
			return 0, fmt.Errorf("read error: %w", ErrFileError)
		}
	}
	book, err := dao.GetBookByHash(database, ctx, hash, false)
	if err != nil {
		if errors.Is(err, dao.ErrDataNotFound) {
			// new book
			log.Logger.Trace().Str("file", fileName).Str("Hash", hash).Msg("Not found in database")
			if !forceUpdate {
				book, err = createBook(basePath, coverPath, fileName, hash)
				if err != nil {
					log.Logger.Error().Str("file", fileName).Str("Hash", hash).Err(err).Msg("File error")
					observer.Error(fileName, "Parse Error")
					return 0, fmt.Errorf("read error: %w", ErrFileError)
				}
			} else {
				book = fileBook
			}
			bookId, err := dao.CreateNewBook(database, ctx, book)
			if err != nil {
				log.Logger.Error().Str("file", fileName).Str("Hash", hash).Err(err).Msg("Book create error")
				observer.Error(fileName, "Database Error")
				return 0, ErrDatabase
			}
			log.Logger.Info().Int("BookId", bookId).Str("file", fileName).Str("Hash", hash).Msg("New book found and created")
			observer.Success(fileName, "Created")
			return int64(bookId), nil

		} else {
			log.Logger.Error().Str("file", fileName).Str("Hash", hash).Err(err).Msg("Database error")
			observer.Error(fileName, "Database Error")
			return 0, ErrDatabase
		}
	}
	log.Logger.Debug().Str("file", fileName).Msg("Old book found")
	relativePath, err := getRelativePath(basePath, fileName)
	if err != nil {
		log.Logger.Error().Str("file", fileName).Str("Hash", hash).Err(err).Msg("Relative path error")
		observer.Error(fileName, "Parse Error")
		return 0, ErrFileError
	}
	relativePath = filepath.ToSlash(relativePath)

	if book.File != relativePath {
		oldPath := book.File
		if !forceUpdate {
			if book.Hascover > 0 {
				if err = os.Rename(path.Join(coverPath, book.File+"."+book.CoverType.String), filepath.Join(coverPath, relativePath+"."+book.CoverType.String)); err != nil {
					log.Logger.Error().Str("cover", relativePath+"."+book.CoverType.String).Err(err).Msg("Rename Error")
					observer.Error(fileName, "Filesystem Error")
					return 0, ErrFileError
				}
			}
			book.File = relativePath
			err = dao.MoveBook(database, ctx, book.Bookid.Int64, relativePath)
			if err != nil {
				log.Logger.Error().Str("file", fileName).Str("Hash", hash).Err(err).Msg("Book Moving Error")
				observer.Error(fileName, "Database Error")
				return 0, ErrDatabase
			}
			log.Logger.Warn().Int64("BookId", book.Bookid.Int64).Str("file", fileName).Str("Old place", oldPath).Msg("Book moved")
			observer.Success(fileName, "Found, moved")
		} else {
			if book.Hascover > 0 {
				err = os.Remove(path.Join(coverPath, book.File+"."+book.CoverType.String))
				if err != nil {
					log.Logger.Error().Str("file", fileName).Str("Hash", hash).Err(err).Msg("Cover delete error")
				}
			}
		}
	} else {
		log.Logger.Info().Int64("BookId", book.Bookid.Int64).Str("file", fileName).Msg("Nothing to do")
	}
	if forceUpdate {
		err = dao.UpdateBookFile(database, ctx, fileBook)
		if err != nil {
			log.Logger.Error().Str("file", fileName).Str("Hash", hash).Err(err).Msg("Book Moving Error")
			observer.Error(fileName, "Database Error")
			return 0, ErrDatabase
		}

	}
	return book.Bookid.Int64, nil
}
func parseADir(observer Observer, basePath string, coverPath string, directory string, forceUpdate bool) ([]int64, error) {
	log.Logger.Debug().Str("Directory", directory).Msg("Enter directory")
	nextDirs := make([]string, 0)
	entries, err := os.ReadDir(directory)
	ret := make([]int64, 0)
	var databaseError error
	databaseError = nil
	if err != nil {
		log.Logger.Error().Str("Directory", directory).Err(err).Msg("Reading directory")
		return ret, ErrFileError
	}
	for _, e := range entries {
		if e.IsDir() {
			nextDirs = append(nextDirs, e.Name())
		} else {
			ext := filepath.Ext(e.Name())
			log.Logger.Debug().Str("file", e.Name()).Str("Ext", ext).Msg("Check File")

			if strings.EqualFold(ext, ".prc") || strings.EqualFold(ext, ".mobi") || strings.EqualFold(ext, ".epub") {
				filename := filepath.Join(directory, e.Name())
				bookId, err := readAFile(observer, basePath, coverPath, filename, forceUpdate)
				if errors.Is(err, ErrDatabase) {
					databaseError = err
				}
				if err == nil {
					ret = append(ret, bookId)
				}
			}
		}
	}
	for _, d := range nextDirs {
		nextPath := filepath.Join(directory, d)
		bookIds, err := parseADir(observer, basePath, coverPath, nextPath, forceUpdate)
		if errors.Is(err, ErrDatabase) {
			databaseError = err
		}
		if err != nil {
			log.Logger.Error().Str("Directory", nextPath).Err(err).Msg("Parsing directory")
		}
		ret = append(ret, bookIds...)
	}
	log.Logger.Debug().Str("Directory", directory).Err(err).Msg("Leave directory")
	return ret, databaseError
}
func Parse(observer Observer, basePath string, coverPath string, forceUpdate bool) error {
	err := observer.StartProcess()
	if err != nil {
		return err
	}
	bookIds, err := parseADir(observer, basePath, coverPath, basePath, forceUpdate)
	if !errors.Is(err, ErrDatabase) {

		slices.Sort(bookIds)
		//log.Logger.Trace().Ints64("Found Ids", bookIds).Msg("Found Books")
		database := db.GetDatabase()
		ctx := context.Background()

		allIDs, err := dao.QueryAllBookId(database, ctx)
		if err == nil {
			log.Logger.Trace().Int("Qty", len(allIDs)).Msg("Books In DB")
			result := utils.Subtract(allIDs, bookIds)
			log.Logger.Trace().Int("Qty", len(result)).Msg("Books not found")
			log.Logger.Info().Ints64("Ids", result).Msg("Books will delete")
			/*
				for _, i := range result {
					err = dao.DeleteBook(database, ctx, i)
					if err != nil {
						log.Logger.Debug().Int64("Book Id", i).Msg("Book deleted")
					}
				}
			*/
		}
	} else {
		return err
	}
	return observer.EndProcess()
}
