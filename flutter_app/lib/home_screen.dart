import 'package:flutter/material.dart';
import 'package:maze_racer/online_game/online_game_screen.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
          child: Center(
        child: ElevatedButton(
            onPressed: () {
              Navigator.of(context).push(
                  MaterialPageRoute(builder: (context) => OnlineGameScreen()));
            },
            child: Text("Play Maze Racer")),
      )),
    );
  }
}
