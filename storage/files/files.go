package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"myFirstTgBot/lib/e"
	"myFirstTgBot/storage"
	"os"
	"path/filepath"
	"time"
)

const defaultPerm = 0774

type Storage struct {
	basePath string
}

func New(basePath string) Storage {

	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() {
		err = e.WrapIfErr("can't save page", err)
	}()
	fPath := filepath.Join(s.basePath, page.UserName)
	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}
	fName, err := fileName(page)
	if err != nil {
		return err
	}
	fPath = filepath.Join(fPath, fName)
	file, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	if err := gob.NewEncoder(file).Encode(page); err != nil {

		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() {
		err = e.WrapIfErr("can't pick random page", err)
	}()
	fPath := filepath.Join(s.basePath, userName)
	files, err := os.ReadDir(fPath)
	if err != nil {

		return nil, err
	}
	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}
	//rand.Seed(time.Now().UnixNano())
	rand.NewSource(time.Now().UnixNano())
	n := rand.Intn(len(files))
	file := files[n]

	return s.decodePage(filepath.Join(fPath, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}
	fpath := filepath.Join(s.basePath, p.UserName, fileName)
	if err := os.Remove(fpath); err != nil {
		msg := fmt.Sprintf("can't remove file %s", fpath)

		return e.Wrap(msg, err)
	}

	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't check if file exists", err)
	}
	fpath := filepath.Join(s.basePath, p.UserName, fileName)
	switch _, err = os.Stat(fpath); {
	case errors.Is(err, os.ErrNotExist):

		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", fpath)

		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func fileName(p *storage.Page) (string, error) {

	return p.Hash()
}

func (s Storage) decodePage(filepath string) (*storage.Page, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	defer func() { _ = f.Close() }()
	var p storage.Page
	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("can't decode page", err)
	}

	return &p, nil
}
