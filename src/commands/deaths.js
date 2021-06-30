const firebase = require('firebase');

module.exports = async (msg) => {
	let kills = [];
	await firebase.firestore().collection('kills').where('serverId', '==', msg.guild.id).get()
		.then((killsQuery) => {
			kills = killsQuery.docs.map(kill => ({
				victim: kill.data().victim,
			}));
		})
		.catch((error) => {
			console.log('Error getting documents: ', error);
		});

	let players = [];
	players = kills.reduce((unique, o) => {
		if(!unique.some(obj => obj === o.victim)) {
			unique.push(o.victim);
		}
		return unique;
	},[]);

	let playerDeaths = [];
	for (const player of players) {
		const playerDeathCount = kills.filter(kill => kill.victim === player).length;
		playerDeaths.push({player: player, killCount: playerDeathCount});
	}
	playerDeaths = playerDeaths.sort((a,b) => b.killCount - a.killCount);


	let killMsg = '**Most Team Deaths\n**';
	for (let i = 0; i < playerDeaths.length; i++) {
		let playerName = '';
		const player = msg.guild.member(playerDeaths[i].player);
		if (player && player.nickname) {
			playerName = player.nickname;
		} else {
			await msg.client.users.fetch(playerDeaths[i].player)
				.then(user => {
					playerName = user.username;
				});
		}

		killMsg += (i + 1) + '. **' + playerName + '** - ' + playerDeaths[i].killCount + ' TDs\n';
	}

	msg.channel.send(killMsg);
};
