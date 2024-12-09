namespace AdventOfCode.Exercises;

public class Exercise9 : IExercise
{
    public int GetDay() => 9;

    private List<int> Parse(string[] input)
    {
        var row = input.First();
        List<int> parsedList = new();

        for (var i = 0; i < row.Count(); i++)
        {
            var rawVal = row[i].ToString();
            int.TryParse(rawVal, out var val);
            parsedList.Add(val);
        }

        return parsedList;
    }

    private long MoveBlocks(List<int> diskMap)
    {
        long checksum = 0;
        var isFile = true;
        var fileId = 0;
        var streamSize = 0;
        var backIdx = diskMap.Count() - 1;
        var backFileId = diskMap.Count() / 2;

        var debugString = "";

        for (var i = 0; i < diskMap.Count(); i++)
        {
            if (diskMap[i] == -1 || diskMap[i+1] == -1 )
            {
                isFile = !isFile;
                continue;
            }

            var fileBlockSize = diskMap[i];

            if (isFile)
            {
                for (var j = 1; j <= fileBlockSize; j++)
                {
                    checksum += fileId * streamSize;
                    streamSize += 1;
                    debugString = debugString + fileId.ToString();
                }
                fileId++;
            }
            else
            {
                if (i > backIdx) continue;
                var fileSizeFromBack = diskMap[backIdx];
                while (fileBlockSize > 0)
                {
                    checksum += backFileId * streamSize;
                    debugString = debugString + backFileId.ToString();
                    fileSizeFromBack -= 1;
                    diskMap[backIdx] = fileSizeFromBack;
                    streamSize += 1;
                    fileBlockSize--;

                    if (fileSizeFromBack == 0)
                    {
                        diskMap[backIdx] = -1;
                        backIdx -= 2;
                        fileSizeFromBack = diskMap[backIdx];
                        backFileId--;
                    }
                }
                diskMap[backIdx] = fileSizeFromBack;
            }
            isFile = !isFile;
        }
        return checksum;
    }

    private List<(int, int)> BuildFileSize(List<int> diskMap)
    {
        List<(int, int)> fileSizes = new();

        var id = 0;
        for (var i = 0; i < diskMap.Count(); i += 2)
        {
            var val = diskMap[i];
            fileSizes.Add((id, val));
            id++;
        }
        fileSizes.Reverse();
        return fileSizes;
    }

    private (int, int, int)? GetSizeToAdd(List<(int, int)> fileSizes, int maxSize)
    {
        for (var i = 0; i < fileSizes.Count(); i++)
        {
            var val = fileSizes[i];
            if (val.Item2 <= maxSize)
            {
                return (val.Item1, val.Item2, i);
            }
        }
        return null;
    }

    private long MoveWholeBlocks(List<int> diskMap, List<(int, int)> fileSizes)
    {
        var isFile = true;
        var fileId = 0;
        var backIdx = diskMap.Count() - 1;
        var backFileId = diskMap.Count() / 2;

        List<int> newDiskMap = [];
        var debugString = "";

        for (var i = 0; i < diskMap.Count(); i++)
        {
            var fileBlockSize = diskMap[i];

            if (isFile)
            {
                if (fileBlockSize < 0)
                {
                    for (var j = 0; j < Math.Abs(fileBlockSize); j++)
                    {
                        newDiskMap.Add(0);
                    }
                }
                else
                {
                    for (var j = 1; j <= fileBlockSize; j++)
                    {
                        newDiskMap.Add(fileId);
                        debugString = debugString + fileId.ToString();
                    }
                    fileSizes.RemoveAt(fileSizes.Count() - 1);
                }
                fileId++;
            }
            else
            {
                var gapSize = diskMap[i];
                while (gapSize != 0) // populate gap
                {
                    // find something with that size
                    var toAdd = GetSizeToAdd(fileSizes, gapSize);
                    if (!toAdd.HasValue) break;

                    (var idToAdd, var sizeToAdd, var fsIdx) = toAdd.Value;

                    // pop last item
                    fileSizes.RemoveAt(fsIdx);

                    // remove item from diskMap
                    var dmIdx = Math.Max((idToAdd * 2), 0);
                    diskMap[dmIdx] = -sizeToAdd;

                    // shrink gapSize
                    gapSize -= sizeToAdd;

                    for (var j = 0; j < sizeToAdd; j++)
                    {
                        newDiskMap.Add(idToAdd);
                        debugString = debugString + idToAdd.ToString();
                    }
                }

                for (var g = 0; g < gapSize; g++)
                {
                    newDiskMap.Add(0);
                }
            }
            isFile = !isFile;
        }

        return newDiskMap.Select((v, i) => (long)(i * v)).Sum();
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        var diskMap = Parse(input);

        var res = MoveBlocks(diskMap);

        return new Result(res, true);
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var diskMap = Parse(input);
        var fileSizes = BuildFileSize(diskMap);
        var res = MoveWholeBlocks(diskMap, fileSizes);

        return new Result(res, false);
    }
}
