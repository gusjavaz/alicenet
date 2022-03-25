import toml from "@iarna/toml";
import fs from "fs";
import { Artifacts } from "hardhat/types";
import { DEFAULT_CONFIG_OUTPUT_DIR, DEPLOYMENT_LIST_FPATH } from "../constants";
import { readDeploymentList } from "./deploymentConfigUtil";
import {
  getDeployGroup,
  getDeployGroupIndex,
  getDeployType,
} from "./deploymentUtil";
export interface ContractDeploymentInfo {
  contract: string;
  index: number;
}
export type DeploymentList = {
  [key: string]: Array<ContractDeploymentInfo>;
};
export interface DeploymentGroupIndexList {
  [key: string]: number[];
}

export async function getDeploymentList(configDirPath?: string) {
  const path =
    configDirPath === undefined
      ? DEFAULT_CONFIG_OUTPUT_DIR + DEPLOYMENT_LIST_FPATH
      : configDirPath + DEPLOYMENT_LIST_FPATH;
  const config: { deploymentList: Array<string> } = await readDeploymentList(
    path
  );
  let deploymentList: Array<string>;
  if (config.deploymentList === undefined) {
    throw new Error(`failed to fetch deployment list at path ${path}`);
  } else {
    deploymentList = config.deploymentList;
  }
  return deploymentList;
}

export async function transformDeploymentList(deploymentlist: DeploymentList) {
  const list: Array<string> = [];
  for (const group in deploymentlist) {
    for (const item of deploymentlist[group]) {
      list.push(item.contract);
    }
  }
  return list;
}

export async function getSortedDeployList(
  contracts: Array<string>,
  artifacts: Artifacts
) {
  const deploymentList: DeploymentList = {};
  for (const contract of contracts) {
    const deployType: string | undefined = await getDeployType(
      contract,
      artifacts
    );
    let group: string | undefined = await getDeployGroup(contract, artifacts);
    let index = -1;
    if (group !== undefined) {
      const indexString: string | undefined = await getDeployGroupIndex(
        contract,
        artifacts
      );
      if (indexString === undefined) {
        throw new Error(
          "If deploy-group-index is specified a deploy-group-index also should be!"
        );
      }
      try {
        index = parseInt(indexString);
      } catch (error) {
        throw new Error(
          `Failed to convert deploy-group-index for contract ${contract}! deploy-group-index should be an integer!`
        );
      }
    } else {
      group = "general";
    }
    if (deployType !== undefined) {
      if (deploymentList[group] === undefined) {
        deploymentList[group] = [];
      }
      deploymentList[group].push({ contract, index });
    }
  }
  for (const key in deploymentList) {
    if (key !== "general") {
      deploymentList[key].sort((contractA, contractB) => {
        return contractA.index - contractB.index;
      });
    }
  }
  return deploymentList;
}

export async function writeDeploymentList(
  list: Array<string>,
  configDirPath?: string
) {
  const path =
    configDirPath === undefined
      ? DEFAULT_CONFIG_OUTPUT_DIR + DEPLOYMENT_LIST_FPATH
      : configDirPath + DEPLOYMENT_LIST_FPATH;
  const config: { deploymentList: Array<string> } = await readDeploymentList(
    path
  );
  config.deploymentList =
    config.deploymentList === undefined ? [] : config.deploymentList;
  config.deploymentList = list;
  const data = toml.stringify(config);
  fs.writeFileSync(path, data);
  const output = fs.readFileSync(path).toString().split("\n");
  output.unshift(
    "# WARNING: THIS IS AN AUTOGENERATED FILE ANY CHANGES WILL BE LOST REFERENCE DOCUMENTATION FOR MORE INFORMATION"
  );
  fs.writeFileSync(path, output.join("\n"));
}
