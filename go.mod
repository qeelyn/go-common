module github.com/qeelyn/go-common

require (
	github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973 // indirect
	github.com/bradfitz/gomemcache v0.0.0-20170208213004-1952afaa557d
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/coreos/etcd v3.3.8+incompatible
	github.com/coreos/go-semver v0.2.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fsnotify/fsnotify v1.4.7 // indirect
	github.com/go-redis/redis v6.10.2+incompatible
	github.com/go-sql-driver/mysql v1.4.0 // indirect
	github.com/gogo/protobuf v1.0.0 // indirect
	github.com/golang/protobuf v1.1.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/hcl v0.0.0-20180404174102-ef8a98b0bbce // indirect
	github.com/jinzhu/gorm v1.9.1
	github.com/jinzhu/inflection v0.0.0-20180308033659-04140366298a // indirect
	github.com/magiconair/properties v1.8.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mitchellh/mapstructure v0.0.0-20180715050151-f15292f7a699 // indirect
	github.com/opentracing/opentracing-go v1.0.2
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pelletier/go-toml v1.2.0 // indirect
	github.com/pkg/errors v0.8.0 // indirect
	github.com/prometheus/client_golang v0.8.0
	github.com/prometheus/client_model v0.0.0-20171117100541-99fa1f4be8e5 // indirect
	github.com/prometheus/common v0.0.0-20180518154759-7600349dcfe1 // indirect
	github.com/prometheus/procfs v0.0.0-20180612222113-7d6f385de8be // indirect
	github.com/spf13/afero v1.1.1 // indirect
	github.com/spf13/cast v1.2.0 // indirect
	github.com/spf13/jwalterweatherman v0.0.0-20180109140146-7c0cea34c8ec // indirect
	github.com/spf13/pflag v1.0.1 // indirect
	github.com/spf13/viper v1.1.0
	github.com/uber/jaeger-client-go v2.14.0+incompatible
	github.com/uber/jaeger-lib v1.5.0
	github.com/ugorji/go v1.1.1 // indirect
	github.com/vmihailenco/msgpack v3.3.2+incompatible
	go.uber.org/atomic v1.3.2 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.8.0
	golang.org/x/net v0.0.0-20180811021610-c39426892332
	golang.org/x/sys v0.0.0-20180810173357-98c5dad5d1a0 // indirect
	golang.org/x/text v0.3.0 // indirect
	google.golang.org/appengine v1.1.0 // indirect
	google.golang.org/genproto v0.0.0-20180808183934-383e8b2c3b9e // indirect
	google.golang.org/grpc v1.14.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0-20170531160350-a96e63847dc3
	gopkg.in/yaml.v2 v2.2.1 // indirect
)

replace (
	golang.org/x/net => github.com/golang/net v0.0.0-20180811021610-c39426892332
	golang.org/x/sys => github.com/golang/sys v0.0.0-20180810173357-98c5dad5d1a0
	golang.org/x/text => github.com/golang/text v0.3.0
	google.golang.org/appengine => github.com/golang/appengine v1.1.0
	google.golang.org/genproto => github.com/google/go-genproto v0.0.0-20180808183934-383e8b2c3b9e
	google.golang.org/grpc => github.com/grpc/grpc-go v1.14.0
)
