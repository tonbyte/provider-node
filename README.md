
# provider-node

This is a temporary implementation of the provider node for the market.tonyte.com project. It will be replaced by a more advanced version soon.

## Server configuration

Perform basic security configuration of your server if desired. You can refer to [this](https://www.informaticar.net/security-hardening-ubuntu-20-04/) for example.

  

## Build TON Storage

Prepare:
```
sudo apt update
sudo apt upgrade
sudo apt install -y build-essential cmake clang openssl libssl-dev zlib1g-dev gperf libreadline-dev ccache libmicrohttpd-dev pkg-config curl

cd ~
mkdir storage
cd storage

git clone --recurse-submodules https://github.com/ton-blockchain/ton.git
cd ton
git checkout testnet
 
cd ..
mkdir ton-build
cd ton-build
```  

If you are using weak hardware, you can increase the swap file size. Refer to [this link](https://docs.ton.org/develop/howto/compile-swap) for instructions.

Build:
```
export CC=/usr/bin/gcc
export CXX=/usr/bin/g++

cmake -DCMAKE_BUILD_TYPE=Release ../ton
cmake --build . -j<number_of_cores_to_compile>
```

## Configure TON Storage

Run:
```
cd ~/storage/ton-build
mkdir storage-db
wget https://ton.org/global-config.json
storage/storage-daemon/storage-daemon -v 3 -C global-config.json -I <YOUR_PUBLIC_IP>:3333 -p 5555 -D storage-db -P -l storage-db/logs.txt &
```

Create the provider contract using:
`storage/storage-daemon/storage-daemon-cli -I 127.0.0.1:5555 -k storage-db/cli-keys/client -p storage-db/cli-keys/server.pub -c "deploy-provider"`

If the above command is successful, you will see the wallet address of your provider contract. Send at least 1 TON to this address. To check if the contract is deployed, you can view the logs with `tail -f storage-db/logs.txt`. Usually, it takes about 10-20 seconds.

Now you can configure your provider with:
`
storage/storage-daemon/storage-daemon-cli -I 127.0.0.1:5555 -k storage-db/cli-keys/client -p storage-db/cli-keys/server.pub -c "set-provider-params --accept 1 --rate 280000 --max-span 65200 --min-file-size 10485760 --max-file-size 5368709120"
`

**Note:** Don't blindly copy-paste this command. You can change the parameters according to your needs. Refer to `storage/storage-daemon/storage-daemon-cli -I 127.0.0.1:5555 -k storage-db/cli-keys/client -p storage-db/cli-keys/server.pub -c "help provider"` for more information.

You can also set the maximum number of contracts and the total size of files with `set-provider-config`.
## Build provider-node
Build:
```
cd ~
git clone https://github.com/tonbyte/provider-node.git
cd provider-node
go build .
```

Edit config.json.
**sp_cli_path** - path to storage-daemon-cli
**storage_db_path** - path to storage-db
**port** - port to listen
**contract_address** - address of your provider contract from previous step
**has_gateway** - set to true if you plan to do next step "Build gateway"

## Build gateway

To use the gateway, we need the latest Node.js and npm. You can install them with:
```
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install -y nodejs
sudo npm install pm2 -g
```

Get gateway:
```
git clone --recursive https://github.com/ton-blockchain/storage-gateway.git
cd storage-gateway
npm install
```

Create environment file using `nano .env` and fill it with:
```
SERVER_PORT=3000
SERVER_HOST="0.0.0.0"
SERVER_HOSTNAME="domain.ton"

TONSTORAGE_BIN="/home/dolfo/storage/ton-build/storage/storage-daemon/storage-daemon-cli"
TONSTORAGE_HOST="127.0.0.1:5555"
TONSTORAGE_DATABASE="/home/dolfo/storage/ton-build/storage-db"
TONSTORAGE_TIMEOUT=5000

SESSION_COOKIE_NAME="sid"
SESSION_COOKIE_PASSWORD="password-should-be-32-characters"
SESSION_COOKIE_ISSECURE=false

GITHUB_AUTH_PASSWORD="password-should-be-32-characters"
GITHUB_AUTH_CLIENTID="authcliendid"
GITHUB_AUTH_CLIENTSECRET="authclientsecret"
GITHUB_AUTH_ISSECURE=false
```
**Note:** If you compiled storage-daemon in a different directory, you should change `TONSTORAGE_BIN` and `TONSTORAGE_DATABASE`. To save the file in `nano`, press Ctrl+X, then Y, and Enter.

Edit `src/config.js` and change `autoload` to `false` and `maxFileSize` to `10737418240` (10GB).
## Configure nginx

Before proceeding with this step, you should have a domain name. If you need assistance, you can contact me on Telegram: [https://dearjohndoe.t.me](https://dearjohndoe.t.me).

Increase the max file size in the nginx configuration:
```
sudo nano /etc/nginx/nginx.conf
```
Find the `http` section and add `client_max_body_size 10G;` inside it. To save the file in `nano`, press Ctrl+X, then Y, and Enter.

Use the template from `scripts/template.conf`. Set your domain name as the `server_name` value.
Edit the nginx config with the template:
```
sudo nano /etc/nginx/sites-available/default
```

Get SSL certificate:
```
sudo apt install -y snapd
sudo snap install core
sudo snap refresh core
sudo snap install --classic certbot
sudo ln -s /snap/bin/certbot /usr/bin/certbot
sudo certbot --nginx
```

## Run everything
Run gateway:
```
cd ~/storage-gateway
npm start
```
**Note:** To check if gateway is working properly call `pm2 logs`.
  
Run provider-node:
```
cd ~/provider-node
./provider-node > logs.txt &
```

## Check if everything is working

**https://<server_name>/v1/provider/status** - should return json with info about your provider.
