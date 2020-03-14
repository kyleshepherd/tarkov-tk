module.exports = async (msg) => {
	let helpMsg = '**Tarkov TK Help** \n \n';
	helpMsg += '`!start` - This will initalise the bot, if not already done \n \n';
	helpMsg += '`!log @Killer @Victim` - This will log a team kill, where the first tagged user is the killer, and the second is the victim \n \n';
	helpMsg += '`!kills` - This will display a scoreboard of the users with the most team kills \n \n';
	helpMsg += '`!deaths` - This will display a scoreboard of the users with the most team deaths \n \n';
	await msg.channel.send(helpMsg);
};