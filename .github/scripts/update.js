(async () => {
    const [,, file] = process.argv;
    const { readFile, writeFile } = await import('fs/promises');
    const path = `${process.cwd()}/configuration/${file}.json`;
    await writeFile(path, JSON.stringify(JSON.parse(await readFile(path)), null, 4), { encoding: 'utf-8' });
})();