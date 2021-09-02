const firebase = require("firebase");

module.exports = async (msg, args) => {
  if (args.length < 1) {
    // stats for whole server
    let stats = [];
    await firebase
      .firestore()
      .collection("kills")
      .where("serverId", "==", msg.channel.guild.id)
      .orderBy("date", "desc")
      .get()
      .then((statsQuery) => {
        stats = statsQuery.docs.map((kill) => ({
          date: kill.data().date,
          killer: kill.data().killer,
          victim: kill.data().victim,
          reason: kill.data().reason,
        }));
      })
      .catch((error) => {
        console.log("Error getting documents: ", error);
      });

    let statMsg = "**Server Kill Stats**\n";
    for (let i = 0; i < stats.length; i++) {
      let killerName = "";
      let victimName = "";

      const killer = msg.guild.member(stats[i].killer);
      if (killer && killer.nickname) {
        killerName = killer.nickname;
      } else {
        await msg.client.users.fetch(stats[i].killer).then((user) => {
          killerName = user.username;
        });
      }

      const victim = msg.guild.member(stats[i].victim);
      if (victim && victim.nickname) {
        victimName = victim.nickname;
      } else {
        await msg.client.users.fetch(stats[i].victim).then((user) => {
          victimName = user.username;
        });
      }

      const date = stats[i].date.toDate();

      statMsg =
        statMsg +
        date.getDate() +
        "/" +
        (date.getMonth() + 1) +
        "/" +
        date.getFullYear() +
        " - Killer: **" +
        killerName +
        "** - Victim: **" +
        victimName +
        "** ";
      if (stats[i].reason !== "") {
        statMsg += '- Reason: "' + stats[i].reason + '"';
      }
      statMsg += "\n\n";
    }
    if (statMsg.length > 2000) {
      const messages = [];
      for (var i = 0; i < statMsg.length; i += 2000) {
        messages.push(statMsg.substr(i, 2000));
      }
      for (let i = 0; i < messages.length; i++) {
        await msg.channel.send(messages[i]);
      }
    } else {
      await msg.channel.send(statMsg);
    }
  } else {
    // stats for single player
    if (msg.mentions.users.size < 1) {
      await msg.channel.send(
        "Make sure you tag a user to see their stats e.g. `!tk @Player`"
      );
    } else {
      const iterator = msg.mentions.users.values();
      const player = iterator.next().value;

      let stats = [];
      await firebase
        .firestore()
        .collection("kills")
        .where("serverId", "==", msg.guild.id)
        .where("killer", "==", player.id)
        .orderBy("date", "desc")
        .get()
        .then((statsQuery) => {
          stats = statsQuery.docs.map((kill) => ({
            date: kill.data().date,
            victim: kill.data().victim,
            reason: kill.data().reason,
          }));
        })
        .catch((error) => {
          console.log("Error getting documents: ", error);
        });
      let killerName = "";
      const killer = msg.guild.member(player.id);
      if (killer && killer.nickname) {
        killerName = killer.nickname;
      } else {
        await msg.client.users.fetch(player.id).then((user) => {
          killerName = user.username;
        });
      }

      if (stats.length === 0) {
        msg.channel.send(killerName + " hasn't team killed anyone...yet");
      } else {
        let statMsg = "**" + killerName + " Team Kills:** \n \n";

        for (let i = 0; i < stats.length; i++) {
          let victimName = "";
          await msg.client.users.fetch(stats[i].victim).then((user) => {
            victimName = user.username;
          });
          const date = stats[i].date.toDate();

          statMsg +=
            date.getDate() +
            "/" +
            (date.getMonth() + 1) +
            "/" +
            date.getFullYear() +
            " - Victim: **" +
            victimName +
            "** ";
          if (stats[i].reason !== "") {
            statMsg += '- Reason: "' + stats[i].reason + '"';
          }
          statMsg += "\n \n";
        }
        if (statMsg.length > 2000) {
          const messages = [];
          for (let i = 0; i < statMsg.length; i += 2000) {
            messages.push(statMsg.substr(i, 2000));
          }
          for (let i = 0; i < messages.length; i++) {
            await msg.channel.send(messages[i]);
          }
        } else {
          await msg.channel.send(statMsg);
        }
      }
    }
  }
};
