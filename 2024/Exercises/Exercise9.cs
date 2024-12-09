using AdventOfCode.Utilities;

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
                    // Console.WriteLine($"Check: { checksum } | Id: { fileId } | FS: {fileBlockSize} | SS: { streamSize } | IsF: {isFile} | it: {i}");
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
                    // Console.WriteLine($"Check: { checksum } | Id: { backFileId } | FS: {fileSizeFromBack} | SS: { streamSize } | IsF: {isFile}");
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
            // Console.WriteLine(string.Join(", ", diskMap));
        }
        // Console.WriteLine(debugString);
        return checksum;
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        var diskMap = Parse(input);

        var res = MoveBlocks(diskMap);

        return new Result(res, false);
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        return new Result();
    }
}
