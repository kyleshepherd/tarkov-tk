module.exports = async (msg) => {
	// var existing = db_checker(msg);
	// existing.then(async function(result) {
	// 	if (result) {
	// 		var players = await get_players(msg);
	// 		for (let i = 0; i < players.length; i++) {
	// 			let playerKills = await get_player_kills(players[i].player_id, msg);
	// 			players[i].kills = playerKills[0].killCount;
	// 		}
	// 		players = sortByKey(players, 'kills');

	// 		let killMsg = '**Most Team Kills \n**';
	// 		for (let i = 0; i < players.length; i++) {
	// 			var playerName = await get_player_name(msg, players[i].player_id);
	// 			killMsg += (i + 1) + '. **' + playerName + '** - ' + players[i].kills + ' TKs \n';
	// 		}
	// 		await msg.channel.send(killMsg);
	// 	} else {
	// 		await msg.channel.send('Tarkov TK has not been set up on this server. Run `!tkstart` to do so.');
	// 	}
	// });

};
