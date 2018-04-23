package util

import (
	"net"
	"reflect"
	"sync"
	"testing"
	"time"

	"golang.org/x/net/icmp"
)

func TestNewPinger(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name    string
		args    args
		want    *Pinger
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := NewPinger(tt.args.addr)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. NewPinger() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. NewPinger() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestPinger_SetIPAddr(t *testing.T) {
	type fields struct {
		Interval    time.Duration
		Timeout     time.Duration
		Count       int
		Debug       bool
		PacketsSent int
		PacketsRecv int
		rtts        []time.Duration
		OnRecv      func(*Packet)
		OnFinish    func(*Statistics)
		done        chan bool
		ipaddr      *net.IPAddr
		addr        string
		ipv4        bool
		source      string
		size        int
		sequence    int
		network     string
	}
	type args struct {
		ipaddr *net.IPAddr
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		p := &Pinger{
			Interval:    tt.fields.Interval,
			Timeout:     tt.fields.Timeout,
			Count:       tt.fields.Count,
			Debug:       tt.fields.Debug,
			PacketsSent: tt.fields.PacketsSent,
			PacketsRecv: tt.fields.PacketsRecv,
			rtts:        tt.fields.rtts,
			OnRecv:      tt.fields.OnRecv,
			OnFinish:    tt.fields.OnFinish,
			done:        tt.fields.done,
			ipaddr:      tt.fields.ipaddr,
			addr:        tt.fields.addr,
			ipv4:        tt.fields.ipv4,
			source:      tt.fields.source,
			size:        tt.fields.size,
			sequence:    tt.fields.sequence,
			network:     tt.fields.network,
		}
		p.SetIPAddr(tt.args.ipaddr)
	}
}

func TestPinger_IPAddr(t *testing.T) {
	type fields struct {
		Interval    time.Duration
		Timeout     time.Duration
		Count       int
		Debug       bool
		PacketsSent int
		PacketsRecv int
		rtts        []time.Duration
		OnRecv      func(*Packet)
		OnFinish    func(*Statistics)
		done        chan bool
		ipaddr      *net.IPAddr
		addr        string
		ipv4        bool
		source      string
		size        int
		sequence    int
		network     string
	}
	tests := []struct {
		name   string
		fields fields
		want   *net.IPAddr
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		p := &Pinger{
			Interval:    tt.fields.Interval,
			Timeout:     tt.fields.Timeout,
			Count:       tt.fields.Count,
			Debug:       tt.fields.Debug,
			PacketsSent: tt.fields.PacketsSent,
			PacketsRecv: tt.fields.PacketsRecv,
			rtts:        tt.fields.rtts,
			OnRecv:      tt.fields.OnRecv,
			OnFinish:    tt.fields.OnFinish,
			done:        tt.fields.done,
			ipaddr:      tt.fields.ipaddr,
			addr:        tt.fields.addr,
			ipv4:        tt.fields.ipv4,
			source:      tt.fields.source,
			size:        tt.fields.size,
			sequence:    tt.fields.sequence,
			network:     tt.fields.network,
		}
		if got := p.IPAddr(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Pinger.IPAddr() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestPinger_SetAddr(t *testing.T) {
	type fields struct {
		Interval    time.Duration
		Timeout     time.Duration
		Count       int
		Debug       bool
		PacketsSent int
		PacketsRecv int
		rtts        []time.Duration
		OnRecv      func(*Packet)
		OnFinish    func(*Statistics)
		done        chan bool
		ipaddr      *net.IPAddr
		addr        string
		ipv4        bool
		source      string
		size        int
		sequence    int
		network     string
	}
	type args struct {
		addr string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		p := &Pinger{
			Interval:    tt.fields.Interval,
			Timeout:     tt.fields.Timeout,
			Count:       tt.fields.Count,
			Debug:       tt.fields.Debug,
			PacketsSent: tt.fields.PacketsSent,
			PacketsRecv: tt.fields.PacketsRecv,
			rtts:        tt.fields.rtts,
			OnRecv:      tt.fields.OnRecv,
			OnFinish:    tt.fields.OnFinish,
			done:        tt.fields.done,
			ipaddr:      tt.fields.ipaddr,
			addr:        tt.fields.addr,
			ipv4:        tt.fields.ipv4,
			source:      tt.fields.source,
			size:        tt.fields.size,
			sequence:    tt.fields.sequence,
			network:     tt.fields.network,
		}
		if err := p.SetAddr(tt.args.addr); (err != nil) != tt.wantErr {
			t.Errorf("%q. Pinger.SetAddr() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestPinger_Addr(t *testing.T) {
	type fields struct {
		Interval    time.Duration
		Timeout     time.Duration
		Count       int
		Debug       bool
		PacketsSent int
		PacketsRecv int
		rtts        []time.Duration
		OnRecv      func(*Packet)
		OnFinish    func(*Statistics)
		done        chan bool
		ipaddr      *net.IPAddr
		addr        string
		ipv4        bool
		source      string
		size        int
		sequence    int
		network     string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		p := &Pinger{
			Interval:    tt.fields.Interval,
			Timeout:     tt.fields.Timeout,
			Count:       tt.fields.Count,
			Debug:       tt.fields.Debug,
			PacketsSent: tt.fields.PacketsSent,
			PacketsRecv: tt.fields.PacketsRecv,
			rtts:        tt.fields.rtts,
			OnRecv:      tt.fields.OnRecv,
			OnFinish:    tt.fields.OnFinish,
			done:        tt.fields.done,
			ipaddr:      tt.fields.ipaddr,
			addr:        tt.fields.addr,
			ipv4:        tt.fields.ipv4,
			source:      tt.fields.source,
			size:        tt.fields.size,
			sequence:    tt.fields.sequence,
			network:     tt.fields.network,
		}
		if got := p.Addr(); got != tt.want {
			t.Errorf("%q. Pinger.Addr() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestPinger_SetPrivileged(t *testing.T) {
	type fields struct {
		Interval    time.Duration
		Timeout     time.Duration
		Count       int
		Debug       bool
		PacketsSent int
		PacketsRecv int
		rtts        []time.Duration
		OnRecv      func(*Packet)
		OnFinish    func(*Statistics)
		done        chan bool
		ipaddr      *net.IPAddr
		addr        string
		ipv4        bool
		source      string
		size        int
		sequence    int
		network     string
	}
	type args struct {
		privileged bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		p := &Pinger{
			Interval:    tt.fields.Interval,
			Timeout:     tt.fields.Timeout,
			Count:       tt.fields.Count,
			Debug:       tt.fields.Debug,
			PacketsSent: tt.fields.PacketsSent,
			PacketsRecv: tt.fields.PacketsRecv,
			rtts:        tt.fields.rtts,
			OnRecv:      tt.fields.OnRecv,
			OnFinish:    tt.fields.OnFinish,
			done:        tt.fields.done,
			ipaddr:      tt.fields.ipaddr,
			addr:        tt.fields.addr,
			ipv4:        tt.fields.ipv4,
			source:      tt.fields.source,
			size:        tt.fields.size,
			sequence:    tt.fields.sequence,
			network:     tt.fields.network,
		}
		p.SetPrivileged(tt.args.privileged)
	}
}

func TestPinger_Privileged(t *testing.T) {
	type fields struct {
		Interval    time.Duration
		Timeout     time.Duration
		Count       int
		Debug       bool
		PacketsSent int
		PacketsRecv int
		rtts        []time.Duration
		OnRecv      func(*Packet)
		OnFinish    func(*Statistics)
		done        chan bool
		ipaddr      *net.IPAddr
		addr        string
		ipv4        bool
		source      string
		size        int
		sequence    int
		network     string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		p := &Pinger{
			Interval:    tt.fields.Interval,
			Timeout:     tt.fields.Timeout,
			Count:       tt.fields.Count,
			Debug:       tt.fields.Debug,
			PacketsSent: tt.fields.PacketsSent,
			PacketsRecv: tt.fields.PacketsRecv,
			rtts:        tt.fields.rtts,
			OnRecv:      tt.fields.OnRecv,
			OnFinish:    tt.fields.OnFinish,
			done:        tt.fields.done,
			ipaddr:      tt.fields.ipaddr,
			addr:        tt.fields.addr,
			ipv4:        tt.fields.ipv4,
			source:      tt.fields.source,
			size:        tt.fields.size,
			sequence:    tt.fields.sequence,
			network:     tt.fields.network,
		}
		if got := p.Privileged(); got != tt.want {
			t.Errorf("%q. Pinger.Privileged() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestPinger_Run(t *testing.T) {
	type fields struct {
		Interval    time.Duration
		Timeout     time.Duration
		Count       int
		Debug       bool
		PacketsSent int
		PacketsRecv int
		rtts        []time.Duration
		OnRecv      func(*Packet)
		OnFinish    func(*Statistics)
		done        chan bool
		ipaddr      *net.IPAddr
		addr        string
		ipv4        bool
		source      string
		size        int
		sequence    int
		network     string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		p := &Pinger{
			Interval:    tt.fields.Interval,
			Timeout:     tt.fields.Timeout,
			Count:       tt.fields.Count,
			Debug:       tt.fields.Debug,
			PacketsSent: tt.fields.PacketsSent,
			PacketsRecv: tt.fields.PacketsRecv,
			rtts:        tt.fields.rtts,
			OnRecv:      tt.fields.OnRecv,
			OnFinish:    tt.fields.OnFinish,
			done:        tt.fields.done,
			ipaddr:      tt.fields.ipaddr,
			addr:        tt.fields.addr,
			ipv4:        tt.fields.ipv4,
			source:      tt.fields.source,
			size:        tt.fields.size,
			sequence:    tt.fields.sequence,
			network:     tt.fields.network,
		}
		p.Run()
	}
}

func TestPinger_run(t *testing.T) {
	type fields struct {
		Interval    time.Duration
		Timeout     time.Duration
		Count       int
		Debug       bool
		PacketsSent int
		PacketsRecv int
		rtts        []time.Duration
		OnRecv      func(*Packet)
		OnFinish    func(*Statistics)
		done        chan bool
		ipaddr      *net.IPAddr
		addr        string
		ipv4        bool
		source      string
		size        int
		sequence    int
		network     string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		p := &Pinger{
			Interval:    tt.fields.Interval,
			Timeout:     tt.fields.Timeout,
			Count:       tt.fields.Count,
			Debug:       tt.fields.Debug,
			PacketsSent: tt.fields.PacketsSent,
			PacketsRecv: tt.fields.PacketsRecv,
			rtts:        tt.fields.rtts,
			OnRecv:      tt.fields.OnRecv,
			OnFinish:    tt.fields.OnFinish,
			done:        tt.fields.done,
			ipaddr:      tt.fields.ipaddr,
			addr:        tt.fields.addr,
			ipv4:        tt.fields.ipv4,
			source:      tt.fields.source,
			size:        tt.fields.size,
			sequence:    tt.fields.sequence,
			network:     tt.fields.network,
		}
		p.run()
	}
}

func TestPinger_finish(t *testing.T) {
	type fields struct {
		Interval    time.Duration
		Timeout     time.Duration
		Count       int
		Debug       bool
		PacketsSent int
		PacketsRecv int
		rtts        []time.Duration
		OnRecv      func(*Packet)
		OnFinish    func(*Statistics)
		done        chan bool
		ipaddr      *net.IPAddr
		addr        string
		ipv4        bool
		source      string
		size        int
		sequence    int
		network     string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		p := &Pinger{
			Interval:    tt.fields.Interval,
			Timeout:     tt.fields.Timeout,
			Count:       tt.fields.Count,
			Debug:       tt.fields.Debug,
			PacketsSent: tt.fields.PacketsSent,
			PacketsRecv: tt.fields.PacketsRecv,
			rtts:        tt.fields.rtts,
			OnRecv:      tt.fields.OnRecv,
			OnFinish:    tt.fields.OnFinish,
			done:        tt.fields.done,
			ipaddr:      tt.fields.ipaddr,
			addr:        tt.fields.addr,
			ipv4:        tt.fields.ipv4,
			source:      tt.fields.source,
			size:        tt.fields.size,
			sequence:    tt.fields.sequence,
			network:     tt.fields.network,
		}
		p.finish()
	}
}

func TestPinger_Statistics(t *testing.T) {
	type fields struct {
		Interval    time.Duration
		Timeout     time.Duration
		Count       int
		Debug       bool
		PacketsSent int
		PacketsRecv int
		rtts        []time.Duration
		OnRecv      func(*Packet)
		OnFinish    func(*Statistics)
		done        chan bool
		ipaddr      *net.IPAddr
		addr        string
		ipv4        bool
		source      string
		size        int
		sequence    int
		network     string
	}
	tests := []struct {
		name   string
		fields fields
		want   *Statistics
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		p := &Pinger{
			Interval:    tt.fields.Interval,
			Timeout:     tt.fields.Timeout,
			Count:       tt.fields.Count,
			Debug:       tt.fields.Debug,
			PacketsSent: tt.fields.PacketsSent,
			PacketsRecv: tt.fields.PacketsRecv,
			rtts:        tt.fields.rtts,
			OnRecv:      tt.fields.OnRecv,
			OnFinish:    tt.fields.OnFinish,
			done:        tt.fields.done,
			ipaddr:      tt.fields.ipaddr,
			addr:        tt.fields.addr,
			ipv4:        tt.fields.ipv4,
			source:      tt.fields.source,
			size:        tt.fields.size,
			sequence:    tt.fields.sequence,
			network:     tt.fields.network,
		}
		if got := p.Statistics(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Pinger.Statistics() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestPinger_recvICMP(t *testing.T) {
	type fields struct {
		Interval    time.Duration
		Timeout     time.Duration
		Count       int
		Debug       bool
		PacketsSent int
		PacketsRecv int
		rtts        []time.Duration
		OnRecv      func(*Packet)
		OnFinish    func(*Statistics)
		done        chan bool
		ipaddr      *net.IPAddr
		addr        string
		ipv4        bool
		source      string
		size        int
		sequence    int
		network     string
	}
	type args struct {
		conn *icmp.PacketConn
		recv chan<- *packet
		wg   *sync.WaitGroup
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		p := &Pinger{
			Interval:    tt.fields.Interval,
			Timeout:     tt.fields.Timeout,
			Count:       tt.fields.Count,
			Debug:       tt.fields.Debug,
			PacketsSent: tt.fields.PacketsSent,
			PacketsRecv: tt.fields.PacketsRecv,
			rtts:        tt.fields.rtts,
			OnRecv:      tt.fields.OnRecv,
			OnFinish:    tt.fields.OnFinish,
			done:        tt.fields.done,
			ipaddr:      tt.fields.ipaddr,
			addr:        tt.fields.addr,
			ipv4:        tt.fields.ipv4,
			source:      tt.fields.source,
			size:        tt.fields.size,
			sequence:    tt.fields.sequence,
			network:     tt.fields.network,
		}
		p.recvICMP(tt.args.conn, tt.args.recv, tt.args.wg)
	}
}

func TestPinger_processPacket(t *testing.T) {
	type fields struct {
		Interval    time.Duration
		Timeout     time.Duration
		Count       int
		Debug       bool
		PacketsSent int
		PacketsRecv int
		rtts        []time.Duration
		OnRecv      func(*Packet)
		OnFinish    func(*Statistics)
		done        chan bool
		ipaddr      *net.IPAddr
		addr        string
		ipv4        bool
		source      string
		size        int
		sequence    int
		network     string
	}
	type args struct {
		recv *packet
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		p := &Pinger{
			Interval:    tt.fields.Interval,
			Timeout:     tt.fields.Timeout,
			Count:       tt.fields.Count,
			Debug:       tt.fields.Debug,
			PacketsSent: tt.fields.PacketsSent,
			PacketsRecv: tt.fields.PacketsRecv,
			rtts:        tt.fields.rtts,
			OnRecv:      tt.fields.OnRecv,
			OnFinish:    tt.fields.OnFinish,
			done:        tt.fields.done,
			ipaddr:      tt.fields.ipaddr,
			addr:        tt.fields.addr,
			ipv4:        tt.fields.ipv4,
			source:      tt.fields.source,
			size:        tt.fields.size,
			sequence:    tt.fields.sequence,
			network:     tt.fields.network,
		}
		if err := p.processPacket(tt.args.recv); (err != nil) != tt.wantErr {
			t.Errorf("%q. Pinger.processPacket() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestPinger_sendICMP(t *testing.T) {
	type fields struct {
		Interval    time.Duration
		Timeout     time.Duration
		Count       int
		Debug       bool
		PacketsSent int
		PacketsRecv int
		rtts        []time.Duration
		OnRecv      func(*Packet)
		OnFinish    func(*Statistics)
		done        chan bool
		ipaddr      *net.IPAddr
		addr        string
		ipv4        bool
		source      string
		size        int
		sequence    int
		network     string
	}
	type args struct {
		conn *icmp.PacketConn
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		p := &Pinger{
			Interval:    tt.fields.Interval,
			Timeout:     tt.fields.Timeout,
			Count:       tt.fields.Count,
			Debug:       tt.fields.Debug,
			PacketsSent: tt.fields.PacketsSent,
			PacketsRecv: tt.fields.PacketsRecv,
			rtts:        tt.fields.rtts,
			OnRecv:      tt.fields.OnRecv,
			OnFinish:    tt.fields.OnFinish,
			done:        tt.fields.done,
			ipaddr:      tt.fields.ipaddr,
			addr:        tt.fields.addr,
			ipv4:        tt.fields.ipv4,
			source:      tt.fields.source,
			size:        tt.fields.size,
			sequence:    tt.fields.sequence,
			network:     tt.fields.network,
		}
		if err := p.sendICMP(tt.args.conn); (err != nil) != tt.wantErr {
			t.Errorf("%q. Pinger.sendICMP() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestPinger_listen(t *testing.T) {
	type fields struct {
		Interval    time.Duration
		Timeout     time.Duration
		Count       int
		Debug       bool
		PacketsSent int
		PacketsRecv int
		rtts        []time.Duration
		OnRecv      func(*Packet)
		OnFinish    func(*Statistics)
		done        chan bool
		ipaddr      *net.IPAddr
		addr        string
		ipv4        bool
		source      string
		size        int
		sequence    int
		network     string
	}
	type args struct {
		netProto string
		source   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *icmp.PacketConn
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		p := &Pinger{
			Interval:    tt.fields.Interval,
			Timeout:     tt.fields.Timeout,
			Count:       tt.fields.Count,
			Debug:       tt.fields.Debug,
			PacketsSent: tt.fields.PacketsSent,
			PacketsRecv: tt.fields.PacketsRecv,
			rtts:        tt.fields.rtts,
			OnRecv:      tt.fields.OnRecv,
			OnFinish:    tt.fields.OnFinish,
			done:        tt.fields.done,
			ipaddr:      tt.fields.ipaddr,
			addr:        tt.fields.addr,
			ipv4:        tt.fields.ipv4,
			source:      tt.fields.source,
			size:        tt.fields.size,
			sequence:    tt.fields.sequence,
			network:     tt.fields.network,
		}
		if got := p.listen(tt.args.netProto, tt.args.source); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Pinger.listen() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_byteSliceOfSize(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := byteSliceOfSize(tt.args.n); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. byteSliceOfSize() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_ipv4Payload(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := ipv4Payload(tt.args.b); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. ipv4Payload() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_bytesToTime(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := bytesToTime(tt.args.b); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. bytesToTime() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_isIPv4(t *testing.T) {
	type args struct {
		ip net.IP
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := isIPv4(tt.args.ip); got != tt.want {
			t.Errorf("%q. isIPv4() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_isIPv6(t *testing.T) {
	type args struct {
		ip net.IP
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := isIPv6(tt.args.ip); got != tt.want {
			t.Errorf("%q. isIPv6() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_timeToBytes(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := timeToBytes(tt.args.t); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. timeToBytes() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
