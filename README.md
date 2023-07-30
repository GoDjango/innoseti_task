# innoseti_task

## Введение

В качестве сети hyperledger-fabric использовал предоставленную разработчиками тестовую сеть.
Описание как ее запустить описано [здесь](https://hyperledger-fabric.readthedocs.io/ru/latest/test_network.html#fabric)
Возможно, сделано все не достаточно качественно, тут можно было бы многое исправить,
но делалось все быстро не вникая во многие детали. 

## Развертывание сети

Для начала нужно было загрузить репозитарий с примерами, инструкция [здесь](https://hyperledger-fabric.readthedocs.io/ru/latest/install.html):

```bash
curl -sSL https://bit.ly/2ysbOFE | bash -s
```

Переходим в каталог, в котором и будем продолжать работу с сетью:

```bash
cd fabric-samples/test-network
```

Объявим переменные окружения:

```bash
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/

# Переменные среды длял организации Org1

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
```

Развернем сеть:

```bash
./network.sh up createChannel -c chan -s couchdb -verbose
./network.sh deployCC -c chan -ccn basic -ccv 1 -ccs 1 -ccp ../../claim/  -ccl go
```

## Примеры команд

Получение списка всех заявок:

```bash
peer chaincode query -C chan -n basic -c '{"Function": "ListClaims", "Args":[]}'
```

Создание заявки:

```bash
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C chan -n basic --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Function": "CreateClaim", "Args":["claim10", "user1", "SOME CONTENT"]}'
```

## Возможные переходы между статусами

Возможно переход между статусами "на рассмотрении" <-> "уточнение данных" не вполне корректен,
но я посчитал иначе. Суть в том, что пользователь может кое-что изменить в статусе уточнения данных,
и отправить заявку снова на исполнение (возможно было бы лучше изменить статус заявки на открытую).

![](/static/src/img/statuses.png)

## Что можно было бы добавить

1. Авторизация пользователей
2. Права пользователям для изменения, исполнения заявок
3. Возможно стоило бы добавить больше полей, к примеру: дата создания, исполнитель заявки и т.д.
4. Возможно стоило вынести файлы смартконтракта в корень директории, но т.к. я работал в корне директории с тестовой сетью, при развертывании получал ошибку превышения размера директории со смартконтрактом
5. Спасибо, что дочитали до конца ход мыслей =)
