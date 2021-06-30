module.exports = async (msg) => {
	// var existing = db_checker(msg);
	// existing.then(async function(result) {
	// 	if (result) {
	// 		var sql = 'TRUNCATE TABLE kills_' + msg.guild.id;
	// 		db.query(sql, async function(err) {
	// 			if (err) throw err;
	// 			var sql = 'TRUNCATE TABLE players_' + msg.guild.id;
	// 			db.query(sql, async function(err) {
	// 				if (err) throw err;
	// 				await msg.channel.send('Your TK server data has been reset');
	// 			});
	// 		});
	// 	} else {
	// 		await msg.channel.send('Tarkov TK has not been set up on this server. Run `!tkstart` to do so.');
	// 	}
	// });
};
