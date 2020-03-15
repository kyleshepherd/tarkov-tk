module.exports = async (msg) => {
	let helpMsg = '**Tarkov TK Help** \n \n';
	helpMsg += '`!tkstart` - This will initalise the bot, if not already done \n \n';
	helpMsg += '`!tklog @Killer @Victim` - This will log a team kill, where the first tagged user is the killer, and the second is the victim.\nYou can also include a reason e.g. `!log @Killer @Victim Killer thought Victim was a Scav` \n \n';
	helpMsg += '`!tkkills` - This will display a scoreboard of the users with the most team kills \n \n';
	helpMsg += '`!tkdeaths` - This will display a scoreboard of the users with the most team deaths \n \n';
	helpMsg += '`!tkinfo` - Some info about the project and the creator, Kyle';
	await msg.channel.send(helpMsg);
};