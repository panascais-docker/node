const path = require("path");
const util = require("util");
const fs = require("fs");
const { default: got } = require("got");

const readFile = util.promisify(fs.readFile);
const writeFile = util.promisify(fs.writeFile);

const url = {
  4: "https://nodejs.org/dist/latest-v4.x/",
  5: "https://nodejs.org/dist/latest-v5.x/",
  6: "https://nodejs.org/dist/latest-v6.x/",
  7: "https://nodejs.org/dist/latest-v7.x/",
  8: "https://nodejs.org/dist/latest-v8.x/",
  9: "https://nodejs.org/dist/latest-v9.x/",
  10: "https://nodejs.org/dist/latest-v10.x/",
  11: "https://nodejs.org/dist/latest-v11.x/",
  12: "https://nodejs.org/dist/latest-v12.x/",
  13: "https://nodejs.org/dist/latest-v13.x/",
  14: "https://nodejs.org/dist/latest-v14.x/",
  15: "https://nodejs.org/dist/latest-v15.x/",
  16: "https://nodejs.org/dist/latest-v16.x/",
  17: "https://nodejs.org/dist/latest-v17.x/",
  18: "https://nodejs.org/dist/latest-v18.x/",
  19: "https://nodejs.org/dist/latest-v19.x/",
  20: "https://nodejs.org/dist/latest-v20.x/",
  21: "https://nodejs.org/dist/latest-v21.x/",
  22: "https://nodejs.org/dist/latest-v22.x/",
  lts: "https://nodejs.org/dist/latest-v20.x/",
  latest: "https://nodejs.org/dist/latest/",
};

const run = async () => {
  const versionsLocation = path.resolve(
    __dirname,
    "../../configuration/versions.json"
  );
  const checksumsLocation = path.resolve(
    __dirname,
    "../../configuration/checksums.json"
  );

  const [
    originalVersionsParsed,
    originalVersions,
    originalChecksumsParsed,
    originalChecksums,
  ] = await Promise.all([
    readFile(versionsLocation, { encoding: "utf-8" }),
    readFile(checksumsLocation, { encoding: "utf-8" }),
  ]).then((result) => {
    return result.flatMap((entry) => {
      const parsed = JSON.parse(entry);
      return [parsed, JSON.stringify(parsed)];
    });
  });

  const versionList = await Promise.all(
    Object.entries(url).map(async ([id, link]) => {
      const regex =
        /node-v(?<major>\d+)\.(?<minor>\d+)\.(?<patch>\d+)\.tar\.xz/gm;

      const { body } = await got(link);

      const { major, minor, patch } = regex.exec(body).groups;

      const version = [major, minor, patch].join(".");

      return [id, version];
    })
  );

  const versions = versionList.reduce((element, [id, version]) => {
    element[id] = version;
    return element;
  }, {});

  const checksumList = await Promise.all(
    versionList.map(async ([id, version]) => {
      const url = `https://unofficial-builds.nodejs.org/download/release/v${version}/SHASUMS256.txt`;
      const regex = new RegExp(
        `(?<hash>[0-9a-f]{64})  node-v${version}-linux-x64-musl.tar.xz`,
        "gm"
      );

      try {
        const { body } = await got(url);

        const { hash } = regex.exec(body).groups;

        return [id, hash];
      } catch {
        return null;
      }
    })
  ).then((list) => list.filter((entry) => entry));

  const checksums = checksumList.reduce((element, [id, hash]) => {
    element[id] = hash;
    return element;
  }, {});

  const versionKeysBefore = Object.keys(originalVersionsParsed).length;
  const versionKeysAfter = Object.keys(versions).length;
  const checksumKeysBefore = Object.keys(originalChecksumsParsed).length;
  const checksumKeysAfter = Object.keys(checksums).length;
  const changed =
    versionKeysBefore === versionKeysAfter &&
    checksumKeysBefore === checksumKeysAfter &&
    !(
      originalVersions === JSON.stringify(versions) &&
      originalChecksums === JSON.stringify(checksums)
    );

  if (changed) {
    await writeFile(versionsLocation, JSON.stringify(versions, null, 4) + "\n");
    await writeFile(
      checksumsLocation,
      JSON.stringify(checksums, null, 4) + "\n"
    );
    return "continue";
  } else {
    return "exit";
  }
};

module.exports = {
  run,
};
