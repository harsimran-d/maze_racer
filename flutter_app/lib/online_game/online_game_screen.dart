import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:maze_racer/maze/maze_controller.dart';
import 'package:maze_racer/maze/player_position_provider.dart';
import 'package:maze_racer/online_game/player_widget.dart';

class OnlineGameScreen extends ConsumerStatefulWidget {
  const OnlineGameScreen({super.key});

  @override
  ConsumerState<OnlineGameScreen> createState() => OnlineGameScreenState();
}

class OnlineGameScreenState extends ConsumerState<OnlineGameScreen> {
  final FocusNode _keyboard = FocusNode();

  @override
  void initState() {
    super.initState();
    _keyboard.requestFocus();
  }

  @override
  void dispose() {
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final maze = ref.watch(mazeProvider);
    return Scaffold(
      body: SafeArea(
          child: maze == null
              ? Center(
                  child: CircularProgressIndicator(),
                )
              : Padding(
                  padding: EdgeInsets.all(50),
                  child: Column(
                    children: [
                      LayoutBuilder(builder: (context, size) {
                        final boxSize = (size.maxWidth > size.maxHeight
                                ? size.maxHeight
                                : size.maxWidth) /
                            maze.size;

                        return AspectRatio(
                          aspectRatio: 1,
                          child: Stack(
                            children: [
                              GridView.builder(
                                  itemCount: maze.size * maze.size,
                                  gridDelegate:
                                      SliverGridDelegateWithFixedCrossAxisCount(
                                    crossAxisCount: maze.size,
                                    childAspectRatio: 1,
                                  ),
                                  itemBuilder: (context, index) {
                                    final row = index % maze.size;
                                    final column = index ~/ maze.size;

                                    return Container(
                                      decoration: BoxDecoration(
                                          color: maze.mazeMap[row][column] == 1
                                              ? Colors.black
                                              : Colors.white),
                                    );
                                  }),
                              PlayerWidget(
                                size: boxSize,
                              )
                            ],
                          ),
                        );
                      }),
                      if (!kIsWeb)
                        ...["UP", "DOWN", "LEFT", "RIGHT"]
                            .map((text) => ElevatedButton(
                                onPressed: () {
                                  ref.read(playersProvider.notifier).move(text);
                                },
                                child: Text(text))),
                      if (kIsWeb)
                        KeyboardListener(
                          focusNode: _keyboard,
                          child: SizedBox.shrink(),
                          onKeyEvent: (value) {
                            switch (value.character?.toLowerCase()) {
                              case 's':
                                ref.read(playersProvider.notifier).move("DOWN");
                                break;
                              case 'a':
                                ref.read(playersProvider.notifier).move("LEFT");
                                break;
                              case 'd':
                                ref
                                    .read(playersProvider.notifier)
                                    .move("RIGHT");
                                break;
                              case 'w':
                                ref.read(playersProvider.notifier).move("UP");
                                break;
                            }
                          },
                        )
                    ],
                  ),
                )),
    );
  }
}
