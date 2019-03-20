package common

import (
	"io/ioutil"
	"os"

	"ic_server/defines"
)

func ReadPage(page string) ([]byte, error) {
  path := defines.ResRoot;
  path += page;
  f, err := os.Open(path);
  if err != nil {
    return nil, err;
  }

  return ioutil.ReadAll(f);
}