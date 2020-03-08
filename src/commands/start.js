var client = require('../db_helper');

module.exports = async (msg) => {
	const tableCheck = 'SELECT 1 from kills LIMIT 1;';
	client.query(tableCheck, async function (err, result) {
		if (result) {
			//Table exists
			await msg.channel.send('Tarkov TK has already been setup on this server. Use !help to see how to use Tarkov TK.');
		} else {
			//Table doesn't exist
			const createTable = 'CREATE TABLE kills (id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY, killer VARCHAR(255), victim VARCHAR(255));';
			client.query(createTable, async function (err) {
				if (err) throw err;
				console.log('Table created');
				await msg.channel.send('Tarkov TK is setup and ready to use. Use !help to see how to use Tarkov TK.');
			});
		}
	});
};