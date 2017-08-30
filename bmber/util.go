package bmber

import "io/ioutil"

func (b *Bmber) GetPageBytes(u string) ([]byte, error) {
	res, err := b.client.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}
