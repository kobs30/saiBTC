Methods:




Create BTC keys

Post Content-Type: application/x-www-form-urlencoded parameters:

method: generateBTC

Return

{"Private":"L4VNWDt966NgQws4Bz78yqJSUWwmoMthnWuJtc1hCDvPLRzhk3BG","Public":"02b8b66ab46eeec6d4fb0540131098561f2bb8a4ded57b9e7be3fa97a08cddb935","Address":"16ifKwXLbdxSScAeAGwnwp6847JLqct6Uj"}


Sign message

Post Content-Type: application/x-www-form-urlencoded parameters:

method: signMessage

p:<private key>

message:<message>

Return

{"message":"Я засыпаю на ходу","signature":"H7gmNU+SEuJ8pCcig/E113UH4U0WpErH9aQzSf/r9novLjgAij+imLkGhnOCqx2HBan59jAkgnUCaCU/EIAHEd0="}




Validate signature

Post Content-Type: application/x-www-form-urlencoded parameters:

method:validateSignature

a:<btc address of signer>

signature:<signature>

message:<signed message>

Return

{"address":"16ifKwXLbdxSScAeAGwnwp6847JLqct6Uj","message":"Я засыпаю на ходу","signature":"valid"}




Create AES key

Post Content-Type: application/x-www-form-urlencoded parameters:

method: createAESkey

Return

732f933995284a85a5d6e47d83cbc22e4b0ba00dfa1ce37533ca7f76a0fa8d26




Encrypt

Post Content-Type: application/x-www-form-urlencoded parameters:

method: encrypt

message:<message>

k:<AES key>

Return

N2JmZTUyNTBkNjlhNGM1ZDZlNDc3MzY2MDZkMzU2M2RiYjIzOWQzM2Y4ZjZjYWQ2MDM2NmVjODU4YmI4ODc4MzlkNDA5Mzg2NGVlNmVhODVmYzYwMDE0YThlOWM1MTFiNDNhNjQ1NzJjNGEzOGRlMzNlMmIwZQ==




Decrypt

Post Content-Type: application/x-www-form-urlencoded parameters:

method: decrypt

cipher:<cipher>

k:<AES key>

Return

Я засыпаю на ходу
