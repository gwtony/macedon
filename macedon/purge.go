package macedon

type PurgeContext struct {
	ips []string
	cmd string
	log *Log
}

func InitPurgeContext(ips string, port string, cmd string, log *Log) (*PurgeContext, error) {
	pc := &PurgeContext{}

	pc.log = log

	//TODO: deal ips to array

	pc.cmd = cmd

	return pc, nil
}

func (pc *PurgeContext) DoPurge(sc *SshContext) (error){
	pc.log.Debug("Do purge")
	return nil
}
