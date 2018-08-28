# Install and running tutorial

 Currently only solo mode is supported.



## 1. Environment Setup

Requirements

- Go-lang >= 1.9
- Docker >= 17.12.0
- Rabbitmq >= 3.7.7
- Dep >= 1.11



1. Install Engine

   ```shell
   go get it-chain/engine
   ```

2. Move to engine folder

   ```shell
   cd $GOPATH/src/github.com/it-chain/engine
   ```

3. Install dependencies

   ```
   dep ensure
   ```

4. Run engine

   **Make sure to run rabbitmq and docker before run it-chain**

   ```shell
   go install it-chain.go
   it-chain
   ```

   ![[tutorial]run](./images/[tutorial]run.png)

5. Check block on http://localhost:4444/blocks

   ![[tutorial]api-blocks](./images/[tutorial]api-blocks.png)

## 2. Installing a icode

Sample icode url:  https://github.com/junbeomlee/learn-icode

1. Fork learn-icode

   ![[tutorial]fork](./images/[tutorial]fork.png)

2. Add deploy key

   ![[tutorial]sshkey](./images/[tutorial]sshkey.png)

3. Deploy learn-icode

   - `it-chain i deploy [learn-icode-url] [ssh-private-key-path]`

     ![[tutorial]deploy](./images/[tutorial]deploy.png)

     ![[tutorial]deploy-result](./images/[tutorial]deploy-result.png)

   - Check docker container

     ```shell
     docker ps
     ```

     ![[tutorial]docker](./images/[tutorial]docker.png)


   - Check icode on http://localhost:4444/icodes

     ![[tutorial]api-icodes](./images/[tutorial]api-icodes.png)



## 3. Invoke a Transaction

1. Invoke initA function using cmd (**initA** function sets key value A to zero)

```
it-chain ivm invoke [IcodeID] initA
```

![[tutorial]invoke](./images/[tutorial]invoke.png)

2. After invoking a initA, query getA (**getA** function get value of key A)

```
it-chain ivm query [IcodeID] getA
```

![[tutorial]query](./images/[tutorial]query.png)

3. If you invoke incA, you can increase the value of A

```shell
it-chain ivm invoke [IcodeID] incA
```

![[tutorial]incA](./images/[tutorial]incA.png)

4. Check block on https://localhost:4444/blocks

