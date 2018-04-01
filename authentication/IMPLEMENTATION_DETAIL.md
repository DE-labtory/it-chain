Crypto <a name="Crypto"></a>
----------------------------

![crypto-implementation-module](../images/crypto-implementation-module.png)

Crypto contains signing and verification related functions that is used for data in blockchain platform and contains key management functions for the keys used in the process. *it-chain* supports `RSA` and `ECDSA` as signature scheme.

-	KeyGenerator

	KeyGenerator generates a key that matches the signature scheme which is selected for signature process.

-	KeyManager

	KeyManager stores generated key, and loads stored key.

-	Signer

	Signer performs data signature.

-	Verifier

	Verifier verifies the signed data.

-	KeyUtils

	KeyUtils performs necessary processing tasks in the process of storing and loading a key such as converting key data to `PEM` file.

-	Key

	Key provides attribute values related to the interface of key data required in the signature or verification process.

<br>

### Signing process of data

![crypto-implementation-seq](../images/crypto-implementation-seq.png) ​

### Author

[@yojkim](https://github.com/yojkim)
