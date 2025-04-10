module github.com/multiversx/mx-chain-sovereign-bridge-go

go 1.20

replace (
	github.com/multiversx/mx-chain-core-go => github.com/multiversx/mx-chain-core-sovereign-go v1.2.25-0.20250410112225-9b4402144b11
	github.com/multiversx/mx-chain-go => github.com/multiversx/mx-chain-sovereign-go v1.8.12-0.20250410144934-a52862ec607a
)

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/gorilla/websocket v1.5.0
	github.com/joho/godotenv v1.5.1
	github.com/multiversx/mx-chain-core-go v1.2.25-0.20250206111825-25fbb1b4851c
	github.com/multiversx/mx-chain-crypto-go v1.2.13-0.20250410124744-d21d37be8e32
	github.com/multiversx/mx-chain-go v1.7.12-0.20240906130631-c565605bfd25
	github.com/multiversx/mx-chain-logger-go v1.0.15
	github.com/multiversx/mx-sdk-go v1.4.7-0.20250410145633-dd8659d5fe1f
	github.com/stretchr/testify v1.10.0
	github.com/urfave/cli v1.22.16
	google.golang.org/grpc v1.60.1
)

require (
	filippo.io/edwards25519 v1.0.0 // indirect
	github.com/bits-and-blooms/bitset v1.10.0 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.2 // indirect
	github.com/btcsuite/btcd/btcutil v1.1.3 // indirect
	github.com/bytedance/sonic v1.11.6 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/consensys/bavard v0.1.13 // indirect
	github.com/consensys/gnark-crypto v0.12.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.5 // indirect
	github.com/crate-crypto/go-kzg-4844 v0.7.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/denisbrodbeck/machineid v1.0.1 // indirect
	github.com/ethereum/c-kzg-4844 v0.4.0 // indirect
	github.com/ethereum/go-ethereum v1.13.15 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.20.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.5-0.20220116011046-fa5810519dcb // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/golang-lru v0.6.0 // indirect
	github.com/holiman/uint256 v1.2.4 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiversx/concurrent-map v0.1.4 // indirect
	github.com/multiversx/mx-chain-communication-go v1.1.1 // indirect
	github.com/multiversx/mx-chain-storage-go v1.0.19 // indirect
	github.com/multiversx/mx-chain-vm-common-go v1.5.17-0.20241119132002-2fa80c5ec516 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pelletier/go-toml v1.9.3 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/supranational/blst v0.3.11 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20220721030215-126854af5e6d // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/tyler-smith/go-bip39 v1.1.0 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	golang.org/x/arch v0.8.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/exp v0.0.0-20231110203233-9a3e6036ecaa // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231002182017-d307bd883b97 // indirect
	google.golang.org/protobuf v1.36.3 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	rsc.io/tmplfunc v0.0.3 // indirect
)
