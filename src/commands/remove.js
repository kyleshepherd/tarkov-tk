const firebase = require('firebase');
const { MessageButton } = require('discord-buttons');

module.exports = async (msg, args) => {
	let kills = [];
	if (args.length < 1) {
		await firebase.firestore().collection('kills').where('serverId', '==', msg.guild.id).orderBy('date', 'desc').limit(5).get()
			.then((killsQuery) => {
				kills = killsQuery.docs.map(kill => ({
					id: kill.id,
					date: kill.data().date,
					killer: kill.data().killer,
					victim: kill.data().victim,
					reason: kill.data().reason,
				}));
			})
			.catch((error) => {
				console.log('Error getting documents: ', error);
			});
	} else {
		if (msg.mentions.users.size < 1) {
			await msg.channel.send('Make sure you tag a user to get their kills e.g. `!tkremove @Player`');
		} else {
			const iterator = msg.mentions.users.values();
			const player = iterator.next().value;

			await firebase.firestore().collection('kills').where('serverId', '==', msg.guild.id).where('killer', '==', player.id).orderBy('date', 'desc').limit(5).get()
				.then((killsQuery) => {
					kills = killsQuery.docs.map(kill => ({
						id: kill.id,
						date: kill.data().date,
						killer: kill.data().killer,
						victim: kill.data().victim,
						reason: kill.data().reason,
					}));
				})
				.catch((error) => {
					console.log('Error getting documents: ', error);
				});
		}
	}

	if (kills.length === 0) {
		await msg.channel.send('No kills logged');
	}

	for (let i = 0; i < kills.length; i++) {
		let killerName = '';
		let victimName = '';

		const killer = msg.guild.member(kills[i].killer);
		if (killer && killer.nickname) {
			killerName = killer.nickname;
		} else {
			await msg.client.users.fetch(kills[i].killer)
				.then(user => {
					killerName = user.username;
				});
		}

		const victim  = msg.guild.member(kills[i].victim);
		if (victim && victim.nickname) {
			victimName = victim.nickname;
		} else {
			await msg.client.users.fetch(kills[i].victim)
				.then(user => {
					victimName = user.username;
				});
		}

		const date = kills[i].date.toDate();

		let button = new MessageButton()
			.setLabel('Remove kill')
			.setStyle('red')
			.setID(kills[i].id);

		let killMsg = date.getDate() + '/' + (date.getMonth() + 1) + '/' + date.getFullYear() + ' - Killer: **' + killerName + '** - Victim: **' + victimName + '** ';
		if (kills[i].reason !== '') {
			killMsg += '- Reason: "' + kills[i].reason + '"';
		}

		await msg.channel.send(killMsg, button);
	}
};
