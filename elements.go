package seek

import "github.com/tebeka/selenium"

type Elements struct {
	err   error
	wd    selenium.WebDriver
	elems []Element
}

func (e Elements) Elements() []Element {
	return e.elems
}

func (e Elements) Valid() (Elements, error) {
	return e, e.err
}

func (e Elements) Error() error {
	return e.err
}

func (e Elements) Len() int {
	return len(e.elems)
}

func (e Elements) FindElementByText(compare func(text string) bool) (elem Element, err error) {
	if e.err != nil {
		err = e.err
		return
	}
	for _, v := range e.elems {
		if v.err != nil {
			err = v.err
			return
		}
		text, _ := v.elem.Text()
		if compare(text) {
			return v, nil
		}
	}
	err = ElementNotFoundErr
	return
}
