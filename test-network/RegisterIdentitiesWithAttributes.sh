fabric-ca-client register --id.name creator1 --id.secret creator1pw --id.type client --id.affiliation org1 --id.attrs 'abac.creator=true:ecert' --tls.certfiles "${PWD}/organizations/fabric-ca/org1/tls-cert.pem"
fabric-ca-client enroll -u https://creator1:creator1pw@localhost:7054 --caname ca-org1 -M "${PWD}/organizations/peerOrganizations/org1.example.com/users/creator1@org1.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/org1/tls-cert.pem"
cp "${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/org1.example.com/users/creator1@org1.example.com/msp/config.yaml"
fabric-ca-client register --id.name creator2 --id.secret creator2pw --id.type client --id.affiliation org1 --id.attrs 'abac.creator=true:' --tls.certfiles "${PWD}/organizations/fabric-ca/org1/tls-cert.pem"
fabric-ca-client enroll -u https://creator2:creator2pw@localhost:7054 --caname ca-org1 --enrollment.attrs "abac.creator" -M "${PWD}/organizations/peerOrganizations/org1.example.com/users/creator2@org1.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/org1/tls-cert.pem"
cp "${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/org1.example.com/users/creator2@org1.example.com/msp/config.yaml"

