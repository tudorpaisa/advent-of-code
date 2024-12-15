namespace AdventOfCode.Exercises;

public class Exercise15 : IExercise
{
    public int GetDay() => 15;

    private class Coords
    {
        public int X { get; set; }
        public int Y { get; set; }
    }

    private List<IMapObject> GetMovableFromItem(List<List<IMapObject?>> map, IMapObject obj, Coords dir)
    {
        List<IMapObject> movable = new();
        for (var i = 0; i < obj.Width; i++)
        {
            var nextX = obj.Position.X + dir.X  + i;
            var nextY = obj.Position.Y + dir.Y;
            Console.WriteLine($"{nextX} {nextY}");
            if (map[nextY][nextX] == null)
            {
                continue;
            }
            var nextItem = map[nextY][nextX];
            if (nextItem.Movable)
            {
                if (nextItem.Position.X == obj.Position.X && nextItem.Position.Y == obj.Position.Y) { continue; }
                if (movable.Count() > 0)
                {
                    if (movable.Last().Position.X != nextItem.Position.X && movable.Last().Position.Y != nextItem.Position.Y)
                    {
                        movable.Add(nextItem);
                    }

                }
                else
                {
                    movable.Add(nextItem);
                }
            }
        }
        Console.WriteLine(movable.Count());

        return movable;
    }

    private IMapObject? GetAdjacentMovable(List<List<IMapObject?>> map, int x, int y)
    {
        if (map[y][x] != null && map[y][x].Movable) return map[y][x];
        return null;
    }

    private interface IMapObject
    {
        bool Movable { get; }
        Coords Position { get; set; }
        int Width { get; }

        Coords MoveBy(List<List<IMapObject?>> map, int x, int y);
        void Clear(List<List<IMapObject?>> map);
        void Insert(List<List<IMapObject?>> map);
    }

    private class Wall : IMapObject
    {
        public int Width => 1;
        public bool Movable => false;
        public Coords Position { get; set; }

        public Coords MoveBy(List<List<IMapObject?>> map, int x, int y) => Position;
        public void Clear(List<List<IMapObject?>> map) {}
        public void Insert(List<List<IMapObject?>> map) {}

    }

    private class Robot : IMapObject
    {
        public int Width => 1;
        public bool Movable => true;
        public Coords Position { get; set; }

        public Coords MoveBy(List<List<IMapObject?>> map, int x, int y)
        {
            var newX = Position.X + x;
            var newY = Position.Y + y;

            if (map[newY][newX] != null) return Position;

            Position.X = newX;
            Position.Y = newY;

            return Position;
        }

        public void Clear(List<List<IMapObject?>> map)
        {
            map[Position.Y][Position.X] = null;
        }
        public void Insert(List<List<IMapObject?>> map)
        {
            map[Position.Y][Position.X] = this;
        }
    }

    private class Box : IMapObject
    {
        public int Width => 1;
        public bool Movable => true;
        public Coords Position { get; set; }

        public Coords MoveBy(List<List<IMapObject?>> map, int x, int y)
        {
            var newX = Position.X + x;
            var newY = Position.Y + y;

            if (map[newY][newX] != null) return Position;

            Position.X = newX;
            Position.Y = newY;

            return Position;
        }

        public void Clear(List<List<IMapObject?>> map)
        {
            map[Position.Y][Position.X] = null;
        }
        public void Insert(List<List<IMapObject?>> map)
        {
            map[Position.Y][Position.X] = this;
        }
    }

    private class WideBox : IMapObject
    {
        public int Width => 2;
        public bool Movable => true;
        public Coords Position { get; set; }

        public Coords MoveBy(List<List<IMapObject?>> map, int x, int y)
        {
            var newX = Position.X + x;
            var newY = Position.Y + y;

            if (map[newY][newX] != null && map[newY][newX + 1] != null) return Position;

            Position.X = newX;
            Position.Y = newY;

            return Position;
        }

        public void Clear(List<List<IMapObject?>> map)
        {
            map[Position.Y][Position.X] = null;
            var adjIt = map[Position.Y][Position.X+1];
            if (adjIt.Position.X != Position.X && adjIt.Position.Y != Position.Y)
            {
                throw new Exception("WTF");
            }
            map[Position.Y][Position.X+1] = null;
        }
        public void Insert(List<List<IMapObject?>> map)
        {
            map[Position.Y][Position.X] = this;
            var adjIt = map[Position.Y][Position.X+1];
            if (adjIt != null)
            {
                throw new Exception($"WTF NOT NULL {adjIt.Position.X} {adjIt.Position.Y} {adjIt.GetType().ToString()}");
            }
            map[Position.Y][Position.X+1] = this;
        }
    }

    private (List<List<IMapObject?>>, List<Coords>, Robot, List<IMapObject>) Parse(string[] input)
    {
        List<List<IMapObject?>> parsed = new();
        List<Coords> directions = new();
        List<IMapObject> boxes = new();
        Robot robot = new();

        bool finishedMap = false;

        for (var i = 0; i < input.Count(); i++)
        {
            if (!finishedMap)
            {
                if (input[i] == "")
                {
                    finishedMap = true;
                    continue;
                }

                List<IMapObject?> parsedRow = new();
                for (var j = 0; j < input[0].Count(); j++)
                {
                    var c = input[i][j];

                    switch (c)
                    {
                        case '#':
                            parsedRow.Add(new Wall { Position = new() { X = j, Y = i } });
                            break;
                        case 'O':
                            var box = new Box { Position = new() { X = j, Y = i } };
                            boxes.Add(box);
                            parsedRow.Add(box);
                            break;
                        case '[':
                            var wideBox = new WideBox { Position = new() { X = j, Y = i } };
                            boxes.Add(wideBox);
                            parsedRow.AddRange([wideBox, wideBox]);
                            // parsedRow.Add(wideBox);
                            // parsedRow.Add(wideBox);
                            j++;
                            break;
                        // case ']':
                        //     break;
                        case '@':
                            robot = new Robot { Position = new() { X = j, Y = i } };
                            parsedRow.Add(robot);
                            break;
                        default:
                            parsedRow.Add(null);
                            break;
                    }
                }
                parsed.Add(parsedRow);
            }
            else
            {
                foreach (var c in input[i])
                {
                    switch (c)
                    {
                        case '<':
                            directions.Add(new() { X = -1, Y = 0 });
                            break;
                        case '^':
                            directions.Add(new() { X = 0, Y = -1 });
                            break;
                        case '>':
                            directions.Add(new() { X = 1, Y = 0 });
                            break;
                        case 'v':
                            directions.Add(new() { X = 0, Y = 1 });
                            break;
                        case 'V':
                            directions.Add(new() { X = 0, Y = 1 });
                            break;
                        default:
                            break;
                    }
                }
            }
        }
        return (parsed, directions, robot, boxes);
    }

    private void PrintMap(List<List<IMapObject?>> map)
    {
        foreach( var row in map )
        {
            for ( var idx = 0; idx < row.Count(); idx++ )
            {
                var i = row[idx];

                if (i == null)
                {
                    Console.Write(" ");
                }
                else if (i.GetType() == typeof(Robot))
                {
                    Console.Write("@");
                }
                else if (i.GetType() == typeof(Box))
                {
                    Console.Write("O");
                }
                else if (i.GetType() == typeof(WideBox) && i.Position.X == idx)
                {
                    Console.Write("[");
                }
                else if (i.GetType() == typeof(WideBox) && i.Position.X == idx - 1)
                {
                    Console.Write("]");
                }
                else if (i.GetType() == typeof(Wall))
                {
                    Console.Write("#");
                }
            }
            Console.Write("\n");
        }
        Console.Write("-------------------\n");
    }

    private string[] ExpandMap(string[] input)
    {
        List<string> expandedMap = new();
        var finishedMap = false;
        foreach (var row in input)
        {
            if (!finishedMap)
            {
                if (row == "")
                {
                    finishedMap = true;
                    expandedMap.Add("");
                    continue;
                }
                var newRow = "";
                foreach (var c in row)
                {
                    switch (c)
                    {
                        case '#':
                            newRow = newRow + "##";
                            break;
                        case '@':
                            newRow = newRow + "@.";
                            break;
                        case '.':
                            newRow = newRow + "..";
                            break;
                        case 'O':
                            newRow = newRow + "[]";
                            break;
                        default:
                            break;
                    }
                }
                expandedMap.Add(newRow);
            }
            else
            {
                expandedMap.Add(row);
            }
        }
        return expandedMap.ToArray();
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        (var map, var directions, var robot, var boxes) = Parse(input);

        foreach (var dir in directions)
        {
            List<IMapObject> itemsToMove = [robot];
            List<IMapObject> moveableChecks = [robot];

            while (moveableChecks.Count() > 0)
            {
                var item = moveableChecks.First();
                moveableChecks.RemoveAt(0);

                var nextX = item.Position.X + dir.X;
                var nextY = item.Position.Y + dir.Y;

                var next = GetAdjacentMovable(map, nextX, nextY);
                if (next != null)
                {
                    moveableChecks.Add(next);
                    itemsToMove = itemsToMove.Prepend(next).ToList();
                }
            }

            foreach (var i in itemsToMove)
            {
                i.Clear(map);
                i.MoveBy(map, dir.X, dir.Y);
                i.Insert(map);
            }
        }

        return new Result(boxes.Select(b => 100 * b.Position.Y + b.Position.X).Sum(), true);
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var expandedMap = ExpandMap(input);
        Console.WriteLine(string.Join("\n", expandedMap));
        (var map, var directions, var robot, var boxes) = Parse(expandedMap);

        foreach (var dir in directions)
        {
            List<IMapObject> itemsToMove = [robot];
            robot.Clear(map);
            List<IMapObject> moveableChecks = [robot];

            while (moveableChecks.Count() > 0)
            {
                var item = moveableChecks.First();
                // Console.WriteLine($"{item.Position.X} {item.Position.Y}");
                moveableChecks.RemoveAt(0);

                var nextX = item.Position.X + dir.X;
                var nextY = item.Position.Y + dir.Y;

                var next = GetMovableFromItem(map, item, dir);
                if (next.Count() > 0)
                {
                    foreach( var x in next )
                    {
                        if (x.Position.X == item.Position.X && x.Position.Y == item.Position.Y) continue;
                        moveableChecks.AddRange(next);
                        itemsToMove = itemsToMove.Prepend(x).ToList();
                        x.Clear(map);
                    }
                }
            }

            foreach (var i in itemsToMove)
            {
                i.MoveBy(map, dir.X, dir.Y);
                i.Insert(map);
            }

            PrintMap(map);
        }

        return new Result();
    }
}
