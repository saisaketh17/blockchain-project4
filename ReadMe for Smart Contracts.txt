Step 1: Go to the directory hyperledger/fabric_Samples/fabcar
Step 2: then ./startfabric.sh run this script this will create mychannel and two organizations
Step 3: Go to javascript and then run this commands
node enrolladmin.js
node registeruser.js
node createProject.js
Which will create project organization and keep them inside blockchain in CouchdB
Step 4: Go to the directory CreateModelTest and run the following commands to run CreateModel Smart Contract
node enrolladmin.js
node createUserWithDeveloperRole.js 
node createUserWithNonDeveloperRole.js
node InvokeUserWithDeveloperRole.js  
node InvokeUserWithOutDeveloperRole.js 
Step 5: Go to the directory ViewModelTest and run the following commands to run QueryModel Smart Contract
node enrolladmin.js
node createOrganizationUser.js. 
node createNonOrganizationUser.js. 
node InvokeViewModelBeingOrganization.js 
node InvokeViewModelBeingNonOrganization.js.
Step 6: Go to the directory UpdateModelTest and run the following commands to run UpdateModel Smart Contract
node enrolladmin.js
node createUserWithDeveloperRole.js 
node createUserWithNonDeveloperRole.js
node InvokeUserWithDeveloperRole.js  
node InvokeUserWithOutDeveloperRole.js

