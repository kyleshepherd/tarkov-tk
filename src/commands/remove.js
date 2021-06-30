module.exports = async (msg) => {
	// var existing = db_checker(msg);
	// existing.then(async function(result) {
	// 	if (result) {
	// 		var sql = 'SELECT id FROM kills_' + msg.guild.id + ' ORDER BY id DESC LIMIT 1;';
	// 		db.query(sql, async function(err, result) {
	// 			if (err) throw err;

	// 			if (result.length != 0) {
	// 				var id = result[0].id;
	// 				var sql = 'DELETE FROM kills_' + msg.guild.id + ' WHERE id = ?;';
	// 				db.query(sql, id, async function(err) {
	// 					if (err) throw err;

	// 					await msg.channel.send('Last logged kill has been removed');
	// 				});
	// 			} else {
	// 				await msg.channel.send('There are no kills to remove');
	// 			}
	// 		});
	// 	} else {
	// 		await msg.channel.send('Tarkov TK has not been set up on this server. Run `!tkstart` to do so.');
	// 	}
	// });
};
