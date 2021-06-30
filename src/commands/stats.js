const firebase = require('firebase');

module.exports = async (msg, args) => {
	if (args.length < 1) {
		// stats for whole server
		firebase.firestore().collection('kills').where('serverId', '==', msg.guild.id).orderBy('date', 'desc').get()
			.then((statsQuery) => {
				let stats = [];
				stats = statsQuery.docs.map(kill => ({
					date: kill.data().date,
					killer: kill.data().killer,
					victim: kill.data().victim,
					reason: kill.data().reason,
				}));

				let statMsg = '**Server Kill Stats**\n';
				for (let i = 0; i < stats.length; i++) {
					const killerName = msg.guild.member(stats[i].killer).nickname;
					const victimName = msg.guild.member(stats[i].victim).nickname;
					const date = stats[i].date.toDate();

					statMsg += date.getDate() + '/' + (date.getMonth() + 1) + '/' + date.getFullYear() + ' - Killer: **' + killerName + '** - Victim: **' + victimName + '** ';
					if (stats[i].reason !== '') {
						statMsg += '- Reason: "' + stats[i].reason + '"';
					}
					statMsg += '\n\n';
				}

				msg.channel.send(statMsg);
			})
			.catch((error) => {
				console.log('Error getting documents: ', error);
			});
	} else {
		// stats for single player
		if (msg.mentions.users.size < 1) {
			await msg.channel.send('Make sure you tag a user to see their stats e.g. `!log @Player`');
		} else {
			const iterator = msg.mentions.users.values();
			const player = iterator.next().value;

			firebase.firestore().collection('kills').where('serverId', '==', msg.guild.id).where('killer', '==', player.id).orderBy('date', 'desc').get()
				.then((statsQuery) => {
					let stats = [];
					stats = statsQuery.docs.map(kill => ({
						date: kill.data().date,
						victim: kill.data().victim,
						reason: kill.data().reason,
					}));

					const killerName = msg.guild.member(player.id).nickname;

					if (stats.length === 0) {
						msg.channel.send(killerName + ' hasn\'t team killed anyone...yet');
					} else {
						let statMsg = '**' + killerName + ' Team Kills:** \n \n';

						for (let i = 0; i < stats.length; i++) {
							const victimName = msg.guild.member(stats[i].victim).nickname;
							const date = stats[i].date.toDate();

							statMsg += date.getDate() + '/' + (date.getMonth() + 1) + '/' + date.getFullYear()  + ' - Victim: **' + victimName + '** ';
							if (stats[i].reason !== '') {
								statMsg += '- Reason: "' + stats[i].reason + '"';
							}
							statMsg += '\n \n';

						}
						msg.channel.send(statMsg);
					}
				})
				.catch((error) => {
					console.log('Error getting documents: ', error);
				});
		}
	}
};
