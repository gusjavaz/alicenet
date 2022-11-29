import fs from "fs";
import { task } from "hardhat/config";
import { HardhatRuntimeEnvironment } from "hardhat/types";
import {
  getDeploymentList,
  getSortedDeployList,
  transformDeploymentList,
} from "./lib/deployment/deploymentListUtil";
import { getAllContracts } from "./lib/deployment/deploymentUtil";

task("generate-immutable-auth-contract", "Generate contracts")
  .addOptionalParam(
    "input",
    "path to the folder containing the deploymentList file",
    undefined
  )
  .addOptionalParam(
    "output",
    "path to save the generated `ImmutableAuth.sol`",
    undefined
  )
  .setAction(async ({ input, output }, hre) => {
    const contractNameSaltMap = [];
    let immutableContract = "";
    let contracts;

    if (input === undefined) {
      const allContracts = await getAllContracts(hre.artifacts);
      const deploymentList = await getSortedDeployList(
        allContracts,
        hre.artifacts
      );
      contracts = await transformDeploymentList(deploymentList);
    } else {
      contracts = await getDeploymentList(input);
    }

    for (let i = 0; i < contracts.length; i++) {
      const fullyQualifiedName = contracts[i];
      const contractName = fullyQualifiedName.split(":")[1];
      const salt = await getSalt(contractName, hre);
      contractNameSaltMap.push([contractName, salt]);
    }

    immutableContract += templateImmutableFactory;
    for (const contractNameSalt of contractNameSaltMap) {
      const name = contractNameSalt[0];
      const salt = contractNameSalt[1];
      const c = new Contract(name, salt);
      immutableContract += templateContract(c);
    }

    if (output === undefined) {
      output = "contracts/utils/";
    } else {
      await checkUserDirPath(output);
    }
    fs.writeFileSync(output + "ImmutableAuth.sol", immutableContract);
  });

export class Contract {
  name: string;
  salt: string;

  constructor(name: string, salt: string) {
    this.name = name;
    this.salt = salt;
  }
}

const templateImmutableFactory = `// This file is auto-generated by hardhat generate-immutable-auth-contract task. DO NOT EDIT.
// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.16;

import "contracts/utils/DeterministicAddress.sol";
import "contracts/interfaces/IAliceNetFactory.sol";

abstract contract ImmutableFactory is DeterministicAddress {
    address private immutable _factory;
    error OnlyFactory(address sender, address expected);

    modifier onlyFactory() {
      if (msg.sender != _factory) {
        revert OnlyFactory(msg.sender, _factory);
      }
      _;
  }

    constructor(address factory_) {
        _factory = factory_;
    }

    function _factoryAddress() internal view returns (address) {
        return _factory;
    }

}

abstract contract ImmutableALCA is ImmutableFactory {
    address private immutable _alca;
    error OnlyALCA(address sender, address expected);

    modifier onlyALCA() {
        if (msg.sender != _alca) {
            revert OnlyALCA(msg.sender, _alca);
        }
        _;
    }

    constructor() {
        _alca = IAliceNetFactory(_factoryAddress()).lookup(_saltForALCA());
    }

    function _alcaAddress() internal view returns (address) {
        return _alca;
    }

    function _saltForALCA() internal pure returns (bytes32) {
        return 0x41546f6b656e0000000000000000000000000000000000000000000000000000;
    }
}

abstract contract ImmutableALCB is ImmutableFactory {
    address private immutable _alcb;
    error OnlyALCB(address sender, address expected);

    modifier onlyALCB() {
        if (msg.sender != _alcb) {
            revert OnlyALCB(msg.sender, _alcb);
        }
        _;
    }

    constructor() {
        _alcb = IAliceNetFactory(_factoryAddress()).lookup(_saltForALCB());
    }

    function _alcbAddress() internal view returns (address) {
        return _alcb;
    }

    function _saltForALCB() internal pure returns (bytes32) {
        return 0x42546f6b656e0000000000000000000000000000000000000000000000000000;
    }
}
`;

function templateContract(contract: Contract): string {
  let mixedContractName: string;
  if (contract.name.startsWith("ETHDKG")) {
    mixedContractName = "ethdkg" + contract.name.slice(6);
  } else {
    mixedContractName =
      contract.name.charAt(0).toLowerCase() + contract.name.slice(1);
  }

  return `
abstract contract Immutable${contract.name} is ImmutableFactory {

    address private immutable _${mixedContractName};
    error Only${contract.name}(address sender, address expected);

    modifier only${contract.name}() {
        if (msg.sender != _${mixedContractName}) {
          revert Only${contract.name}(msg.sender, _${mixedContractName});
        }
        _;
    }

    constructor() {
        _${mixedContractName} = getMetamorphicContractAddress(${contract.salt}, _factoryAddress());
    }

    function _${mixedContractName}Address() internal view returns(address) {
        return _${mixedContractName};
    }

    function _saltFor${contract.name}() internal pure returns(bytes32) {
        return ${contract.salt};
    }
}
    `;
}

function extractPath(qualifiedName: string) {
  return qualifiedName.split(":")[0];
}

async function checkUserDirPath(path: string) {
  if (path !== undefined) {
    if (!fs.existsSync(path)) {
      console.log(
        "Creating Folder at" + path + " since it didn't exist before!"
      );
      fs.mkdirSync(path);
    }
    if (fs.statSync(path).isFile()) {
      throw new Error("outputFolder path should be to a directory not a file");
    }
  }
}

async function getFullyQualifiedName(
  contractName: string,
  hre: HardhatRuntimeEnvironment
) {
  const artifactPaths = await hre.artifacts.getAllFullyQualifiedNames();
  for (let i = 0; i < artifactPaths.length; i++) {
    if (
      artifactPaths[i].split(":").length > 0 &&
      artifactPaths[i].split(":")[1] === contractName
    ) {
      return String(artifactPaths[i]);
    }
  }
}

async function getSalt(
  contractName: string,
  hre: HardhatRuntimeEnvironment
): Promise<string> {
  const qualifiedName: any = await getFullyQualifiedName(contractName, hre);
  const buildInfo = await hre.artifacts.getBuildInfo(qualifiedName);
  let contractOutput: any;
  let devdoc: any;
  let salt: string = "";
  let saltType: string = "";
  if (buildInfo !== undefined) {
    const path = extractPath(qualifiedName);
    contractOutput = buildInfo.output.contracts[path][contractName];
    devdoc = contractOutput.devdoc;
    salt = devdoc["custom:salt"];
    saltType = devdoc["custom:salt-type"];
  } else {
    console.error("missing salt");
  }
  if (saltType === undefined) {
    salt = hre.ethers.utils.formatBytes32String(salt);
  } else {
    salt = hre.ethers.utils.keccak256(
      hre.ethers.utils
        .keccak256(hre.ethers.utils.formatBytes32String(salt))
        .concat(
          hre.ethers.utils
            .keccak256(hre.ethers.utils.formatBytes32String(saltType))
            .slice(2)
        )
    );
  }
  return salt;
}
