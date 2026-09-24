package desktopruntime

import (
 "bufio"
 "bytes"
 "encoding/json"
 "errors"
 "io"
 "net/url"
 "strings"
)

const RuntimeHostFrameLimit = 16 * 1024
const RuntimeHostReplaceOrigin = "replace_core_origin"

// These frames are private to one inherited anonymous-pipe pair. They are not
// registered on MCP, desktop IPC, a socket, or a named pipe.
type RuntimeHostRequest struct {
 Op string `json:"op"`
 RequestID uint64 `json:"request_id"`
 Epoch string `json:"host_epoch"`
 Generation string `json:"generation"`
 Origin string `json:"origin"`
}
type RuntimeHostAck struct {
 RequestID uint64 `json:"request_id"`
 Epoch string `json:"host_epoch"`
 Generation string `json:"generation"`
 CorePID uint32 `json:"core_pid"`
 OK bool `json:"ok"`
 Error string `json:"error,omitempty"`
}
func ReadRuntimeHostFrame(reader *bufio.Reader, value any) error {
 data,err := reader.ReadSlice('\n')
 if err != nil { return err }
 if len(data)>RuntimeHostFrameLimit { return errors.New("runtime host frame exceeds limit") }
 decoder:=json.NewDecoder(bytes.NewReader(data)); decoder.DisallowUnknownFields()
 if err:=decoder.Decode(value); err!=nil { return err }
 if err:=decoder.Decode(new(any)); err!=io.EOF { return errors.New("runtime host frame has trailing data") }
 return nil
}
func WriteRuntimeHostFrame(writer io.Writer, value any) error {
 data,err:=json.Marshal(value); if err!=nil { return err }; data=append(data,'\n')
 if len(data)>RuntimeHostFrameLimit { return errors.New("runtime host frame exceeds limit") }
 n,err:=writer.Write(data); if err==nil && n!=len(data) { err=io.ErrShortWrite }; return err
}
func (request RuntimeHostRequest) Validate(epoch,generation string,next uint64) error {
 if request.Op!=RuntimeHostReplaceOrigin || request.RequestID!=next || next==0 || request.Epoch=="" || request.Epoch!=epoch || request.Generation=="" || request.Generation!=generation {
  return errors.New("unknown, duplicate or stale runtime host request")
 }
 return validateQuickRuntimeOrigin(request.Origin)
}
func validateQuickRuntimeOrigin(origin string) error {
 if origin=="" { return nil } // Invalidate a lost Quick origin without touching Named credentials.
 parsed,err:=url.Parse(origin)
 if err!=nil || parsed.Scheme!="https" || parsed.User!=nil || parsed.Port()!="" || parsed.Path!="" || parsed.RawQuery!="" || parsed.ForceQuery || parsed.Fragment!="" || !strings.HasSuffix(parsed.Host,".trycloudflare.com") || strings.ContainsAny(parsed.Host,"\\% \t\r\n") || parsed.Host!=strings.ToLower(parsed.Host) {
  return errors.New("runtime host requires a canonical Quick HTTPS origin")
 }
 label:=strings.TrimSuffix(parsed.Host,".trycloudflare.com")
 if label=="" || strings.Contains(label,".") { return errors.New("invalid Quick hostname") }
 for _,c:=range label { if c!='-' && (c<'a'||c>'z') && (c<'0'||c>'9') { return errors.New("invalid Quick hostname") } }
 return nil
}
