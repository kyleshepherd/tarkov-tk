module.exports = async (msg) => {
	let helpMsg = '**Tarkov TK Help** \n \n';
	helpMsg += '`!tklog @Killer @Victim` - This will log a team kill, where the first tagged user is the killer, and the second is the victim.\nYou can also include a reason e.g. `!tklog @Killer @Victim Killer thought Victim was a Scav` \n \n';
	helpMsg += '`!tkremove` - This will display the last 5 logged kills and allow you to remove any of them via a button click\nYou can also tag a player to get just their TKs to remove \n \n';
	helpMsg += '`!tkreset` - This command will reset the TK server data for your channel. **THIS WILL DELETE ALL TK LOGS** \n \n';
	helpMsg += '`!tk` - This command will print all of the team kills logged for your server\n \n';
	helpMsg += '`!tk @Player` - This command will print all of the team kills for a specific player \n \n';
	helpMsg += '`!tkkills` - This will display a scoreboard of the users with the most team kills \n \n';
	helpMsg += '`!tkdeaths` - This will display a scoreboard of the users with the most team deaths \n \n';
	helpMsg += '`!tkinfo` - Some info about the project and the creator, Kyle';
	await msg.channel.send(helpMsg);
};
