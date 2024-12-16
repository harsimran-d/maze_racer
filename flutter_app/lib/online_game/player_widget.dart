import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:maze_racer/maze/player_position_provider.dart';

class PlayerWidget extends ConsumerWidget {
  const PlayerWidget({super.key, required this.size});
  final double size;
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final playerPositionsMap = ref.watch(playersProvider);
    return Stack(
      children: playerPositionsMap.entries.map<Widget>((entry) {
        return Positioned(
          top: entry.value.column * size,
          left: entry.value.row * size,
          child: Character(size: size),
        );
      }).toList(),
    );
  }
}

class Character extends StatelessWidget {
  const Character({
    super.key,
    required this.size,
  });

  final double size;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: size,
      height: size,
      child: Container(
        margin: EdgeInsets.all(size * 0.1),
        decoration: BoxDecoration(
          color: Colors.red,
          borderRadius: BorderRadius.circular(
            (size * 0.3),
          ),
        ),
      ),
    );
  }
}
