package dfuse

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"
)

func generateReqId(custom string) string {
	hash := fmt.Sprintf("%s %s %d %d", custom, time.Now(), rand.Uint32(), rand.Uint32())
	buf := bytes.NewBuffer(nil)
	sum := md5.Sum([]byte(hash))
	encoder := base64.NewEncoder(base64.URLEncoding, buf)
	_, _ = encoder.Write(sum[:])
	_ = encoder.Close()
	return buf.String()[:20]
}
