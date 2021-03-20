package prelude

func Make(target string) (*ExecRet, error) {
	return Exec("make", target)
}
