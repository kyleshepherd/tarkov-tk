module.exports = async (msg) => {
	let helpMsg = '**Tarkov TK Help** \n \n';
	helpMsg += '`!tkstart` - This will initalise the bot, if not already done \n \n';
	helpMsg += '`!tklog @Killer @Victim` - This will log a team kill, where the first tagged user is the killer, and the second is the victim \n \n';
	helpMsg += '`!tkkills` - This will display a scoreboard of the users with the most team kills \n \n';
	helpMsg += '`!tkdeaths` - This will display a scoreboard of the users with the most team deaths \n \n';
	await msg.channel.send(helpMsg);
};