var client = require('../db_helper');

module.exports = async (msg, args) => {
	if (args.length < 2) {
		await msg.channel.send('Make sure you tag 2 users in !log \n e.g. `!log @Killer @Victim`');
	} else {
		if (msg.mentions.users.size < 2) {
			await msg.channel.send('Make sure you tag 2 users in !log \n e.g. `!log @Killer @Victim`');
		} else {
			const iterator = msg.mentions.users.values();

			const killer = iterator.next().value;
			const victim = iterator.next().value;

			const killLog = 'INSERT INTO kills (killer, victim) VALUES ("' + killer.id + '", "' + victim.id + '");';

			client.query(killLog, async function (err) {
				if (err) throw err;
				await msg.channel.send('Kill by ' + killer.username + ' on ' + victim.username + ' logged.');
			});
		}
	}
};